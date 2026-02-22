import logging
from docling.document_converter import DocumentConverter, PdfFormatOption, PowerpointFormatOption, WordFormatOption
from docling.datamodel.base_models import InputFormat
from docling_core.transforms.chunker.base import BaseChunk
from docling_core.transforms.chunker.hybrid_chunker import HybridChunker
from docling_core.transforms.chunker.tokenizer.huggingface import HuggingFaceTokenizer
from transformers import AutoTokenizer, AutoModel
import torch
import torch.nn.functional as F

from utils import collect_chunk_meta

EMBED_MODEL_ID = "sentence-transformers/all-MiniLM-L6-v2"
MAX_TOKENS = 250
BATCH_SIZE = 10
hf_tokenizer = AutoTokenizer.from_pretrained(EMBED_MODEL_ID)
embedding_model = AutoModel.from_pretrained(EMBED_MODEL_ID)
embedding_model.eval()

# finds text in pdfs perfectly (docling parse / pymupdf)
# finds good quality images with printed text, easyOCR / rapidOCR / tesserocr works
# TODO: for embedded images, setup VLM or SuryaORM which uses GPU
doc_converter = DocumentConverter(
    allowed_formats=[InputFormat.PDF, InputFormat.DOCX, InputFormat.PPTX],
    format_options={
        InputFormat.PDF: PdfFormatOption(),
        InputFormat.DOCX: WordFormatOption(),
        InputFormat.PPTX: PowerpointFormatOption()
    }
)

def process_document(data, logger):
    result = doc_converter.convert(data["path"])
    chunks = chunkify(result.document)
    embeddings = embed_chunks(chunks)
    for c in chunks:
        print("\n\n###############")
        print(c.text)
        print(collect_chunk_meta(c))
    print(len(embeddings[0]))

def chunkify(doc):
    tokenizer = HuggingFaceTokenizer(
        tokenizer=hf_tokenizer,
        max_tokens=250
    )
    chunker = HybridChunker(tokenizer=tokenizer)
    return list(chunker.chunk(doc))

def mean_pooling(model_output, attention_mask):
    token_embeddings = model_output[0] #First element of model_output contains all token embeddings
    input_mask_expanded = attention_mask.unsqueeze(-1).expand(token_embeddings.size()).float()
    return torch.sum(token_embeddings * input_mask_expanded, 1) / torch.clamp(input_mask_expanded.sum(1), min=1e-9)

def hf_embed_query(texts: list[str]) -> list[list[float]]:
    encoded = hf_tokenizer(
        texts,
        padding=True,
        truncation=True,
        max_length=256,
        return_tensors="pt"
    )

    with torch.no_grad():
        model_output = embedding_model(**encoded)

    sentence_embeddings = mean_pooling(
        model_output,
        encoded["attention_mask"]
    )

    sentence_embeddings = F.normalize(sentence_embeddings, p=2, dim=1)
    return sentence_embeddings.tolist()

def embed_chunks(chunks: list[BaseChunk]):
    embeddings = []

    for i in range(0, len(chunks), BATCH_SIZE):
        batch = chunks[i:i + BATCH_SIZE]
        texts = [c.text for c in batch if c.text.strip()]
        if not texts:
            continue

        emb = hf_embed_query(texts)
        embeddings.extend(emb)

    return embeddings

if __name__ == "__main__":
    logger = logging.getLogger(__name__)
    process_document({
        'path': "./test_assets/t_digital.pdf"
    }, logger)
