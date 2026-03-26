-- Mengaktifkan ekstensi gen_random_uuid() jika belum ada
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Table: pseo_articles
CREATE TABLE IF NOT EXISTS public.pseo_articles (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    slug TEXT UNIQUE NOT NULL,
    title TEXT NOT NULL,
    content_html TEXT NOT NULL,
    meta_desc TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

-- Indexing untuk optimasi baca per slug pada SSG Next.js
CREATE INDEX IF NOT EXISTS pseo_articles_slug_idx ON public.pseo_articles (slug);

-- Table: b2b_leads
CREATE TABLE IF NOT EXISTS public.b2b_leads (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    business_name TEXT NOT NULL,
    website_status TEXT NOT NULL,
    niche TEXT NOT NULL,
    ai_email_draft TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);
