import redis

REDIS_URL = "redis://localhost:6379"
STREAM = "ingestion_queue"
GROUP = "workers"
CONSUMER = "worker-1"
MIN_IDLE_TIME = 10000

class Stream:
    def __init__(self, logger):
        logger.info("Connecting to Redis")

        self.r = redis.Redis.from_url(REDIS_URL, decode_responses=True)
        self.logger = logger

        try:
            logger.info("Creating redis stream %s", STREAM)
            self.r.xgroup_create(STREAM, GROUP, id="0", mkstream=True)
        except redis.ResponseError:
            pass

    def auto_claim(self):
        v = self.r.xautoclaim(STREAM, GROUP, CONSUMER, MIN_IDLE_TIME, start_id="0-0", count=1)
        _, entries, *_ = v
        if len(entries) > 0:
            return entries[0]
        return None

    def read_group(self):
        v = self.r.xreadgroup(
            GROUP,
            CONSUMER,
            {STREAM: ">"},
            count=1,
            block=5000,
        )
        if v:
            _, entries = v[0]
            entry = entries[0]
            return entry if entry else None
        return None

    def ack(self, msg_id):
        self.r.xack(STREAM, GROUP, msg_id)

    def get_delivery_count(self, msg_id):
        pel_data = self.r.xpending_range(STREAM, GROUP, msg_id, msg_id, 1)[0]
        if not pel_data:
            return None

        delivery_count = pel_data['times_delivered']
        return delivery_count
