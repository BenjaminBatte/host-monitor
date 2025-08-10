
CREATE TABLE IF NOT EXISTS public.checks (
  id           BIGSERIAL PRIMARY KEY,
  host         TEXT        NOT NULL CHECK (length(trim(host)) > 0),
  up           BOOLEAN     NOT NULL,
  latency_ms   INTEGER, 
  packet_loss  REAL        NOT NULL DEFAULT 0 CHECK (packet_loss >= 0 AND packet_loss <= 100),
  checked_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_checks_host_time
  ON public.checks (host, checked_at DESC);


CREATE INDEX IF NOT EXISTS idx_checks_down_time
  ON public.checks (checked_at DESC)
  WHERE NOT up;



CREATE OR REPLACE VIEW public.latest_checks AS
SELECT DISTINCT ON (host)
  id, host, up, latency_ms, packet_loss, checked_at
FROM public.checks
ORDER BY host, checked_at DESC;
