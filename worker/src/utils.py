def collect_chunk_meta(chunk):
    page_number = None
    content_layer = None
    kind = "text"   # default

    doc_items = getattr(chunk.meta, "doc_items", []) or []

    for item in doc_items:
        # ---- page number (first one wins) ----
        if page_number is None:
            prov = getattr(item, "prov", None)
            if prov:
                first = prov[0]
                page_number = getattr(first, "page_no", None)

        # ---- content layer (BODY dominates) ----
        cl = getattr(item, "content_layer", None)
        if cl:
            cl_name = getattr(cl, "name", None) or str(cl)
            if content_layer is None:
                content_layer = cl_name
            if cl_name == "BODY":
                content_layer = "BODY"

        # ---- kind (TABLE dominates) ----
        lbl = getattr(item, "label", None)
        if lbl:
            lbl_name = getattr(lbl, "name", None) or str(lbl)
            if lbl_name == "TABLE":
                kind = "table"

        # fast exit if we already have everything dominant
        if page_number is not None and content_layer == "BODY" and kind == "table":
            break

    return {
        "page_number": page_number,
        "headings": getattr(chunk.meta, "headings", []),
        "content_layer": content_layer,
        "kind": kind,
    }
