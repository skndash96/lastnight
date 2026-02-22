from process import process_document
import logging
from stream import Stream
from repo import Repo

def run():
    logger = logging.getLogger(__name__)

    repo = Repo(logger)
    rs = Stream(logger)

    logger.info("Starting Loop")

    while True:
        entry = None
        is_autoclaimed = True

        entry = rs.auto_claim()
        if entry:
            logger.info("Autoclaimed entry %s", entry)
        else:
            is_autoclaimed = False
            entry = rs.read_group()
            if entry:
                logger.info("Read entry %s", entry)
            else:
                logger.info("No entry read")
                continue

        msg_id, msg = entry

        try:
            j = repo.get_document_info(msg)
            process_document(j, logger)
            repo.post_result(msg, True)
            rs.ack(msg_id)
        except Exception as e:
            if is_autoclaimed:
                delivery_count = rs.get_delivery_count(msg_id)
                if delivery_count is None:
                    continue

                logger.info("Document processing failed after %d attempts %s", delivery_count, msg_id, exc_info=True)

                if delivery_count > 2:
                    logger.error("Document processing failed after multiple attempts")
                    repo.post_result(msg, False)
                    rs.ack(msg_id)
            else:
                logger.error("Document processing failed", exc_info=True)

if __name__ == "__main__":
    run()
