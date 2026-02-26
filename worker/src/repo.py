import psycopg2
from psycopg2.extensions import connection

DB_URL = "postgres://postgres:postgres@localhost:5432/postgres"
BUCKET_NAME = "lastnight"

class Repo:
    def __init__(self, logger):
        conn = psycopg2.connect(DB_URL)
        cur = conn.cursor()
        cur.execute("SELECT version()")
        db_version = cur.fetchone()
        logger.info(f"Connected to PostgreSQL {db_version}")
        self.conn = conn
        self.logger = logger

    def get_document_info(self, msg):
        return msg

    def post_result(self, msg, ok: bool):
        return
