-- EstateSaleFinder.ai Initial Database Schema
-- This migration creates the core tables for the estate sale marketplace

-- Users table (Firebase auth + local profile)
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    firebase_uid TEXT UNIQUE NOT NULL,
    email TEXT NOT NULL,
    user_type VARCHAR(50) DEFAULT 'buyer', -- 'buyer', 'seller', 'professional'
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_users_firebase_uid ON users(firebase_uid);
CREATE INDEX idx_users_email ON users(email);

-- User profiles (extended information)
CREATE TABLE IF NOT EXISTS user_profiles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    phone VARCHAR(20),
    company_name VARCHAR(200), -- for professionals
    bio TEXT,
    avatar_url TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_user_profiles_user_id ON user_profiles(user_id);

-- Estate sales (main listings)
CREATE TABLE IF NOT EXISTS estate_sales (
    id SERIAL PRIMARY KEY,
    seller_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    sale_type VARCHAR(50) DEFAULT 'estate_sale', -- 'estate_sale', 'auction', 'moving_sale'
    status VARCHAR(50) DEFAULT 'draft', -- 'draft', 'published', 'completed', 'cancelled'

    -- Location
    address_line1 VARCHAR(255),
    address_line2 VARCHAR(255),
    city VARCHAR(100) NOT NULL,
    state VARCHAR(50) NOT NULL,
    zip_code VARCHAR(20),
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),

    -- Sale dates/times
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    sale_hours TEXT, -- e.g., "Fri 9am-5pm, Sat 9am-3pm"

    -- Pricing (for posting fees)
    listing_tier VARCHAR(50) DEFAULT 'basic', -- 'basic', 'featured', 'premium'
    payment_status VARCHAR(50) DEFAULT 'unpaid', -- 'unpaid', 'paid', 'refunded'
    amount_paid DECIMAL(10, 2) DEFAULT 0.00,

    -- Metadata
    view_count INTEGER DEFAULT 0,
    featured BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_estate_sales_seller_id ON estate_sales(seller_id);
CREATE INDEX idx_estate_sales_city ON estate_sales(city);
CREATE INDEX idx_estate_sales_state ON estate_sales(state);
CREATE INDEX idx_estate_sales_zip_code ON estate_sales(zip_code);
CREATE INDEX idx_estate_sales_status ON estate_sales(status);
CREATE INDEX idx_estate_sales_sale_type ON estate_sales(sale_type);
CREATE INDEX idx_estate_sales_start_date ON estate_sales(start_date);
CREATE INDEX idx_estate_sales_featured ON estate_sales(featured);

-- Sale images
CREATE TABLE IF NOT EXISTS sale_images (
    id SERIAL PRIMARY KEY,
    sale_id INTEGER REFERENCES estate_sales(id) ON DELETE CASCADE,
    image_url TEXT NOT NULL,
    thumbnail_url TEXT,
    is_primary BOOLEAN DEFAULT FALSE,
    display_order INTEGER DEFAULT 0,
    uploaded_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_sale_images_sale_id ON sale_images(sale_id);
CREATE INDEX idx_sale_images_is_primary ON sale_images(is_primary);

-- Sale items (optional granular inventory within a sale)
CREATE TABLE IF NOT EXISTS sale_items (
    id SERIAL PRIMARY KEY,
    sale_id INTEGER REFERENCES estate_sales(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100), -- 'furniture', 'jewelry', 'antiques', 'electronics', etc.
    estimated_price DECIMAL(10, 2),
    image_url TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_sale_items_sale_id ON sale_items(sale_id);
CREATE INDEX idx_sale_items_category ON sale_items(category);

-- Saved sales (buyer favorites)
CREATE TABLE IF NOT EXISTS saved_sales (
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    sale_id INTEGER REFERENCES estate_sales(id) ON DELETE CASCADE,
    saved_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, sale_id)
);

CREATE INDEX idx_saved_sales_user_id ON saved_sales(user_id);
CREATE INDEX idx_saved_sales_sale_id ON saved_sales(sale_id);

-- Professional directory (auctioneers, appraisers, etc.)
CREATE TABLE IF NOT EXISTS professionals (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    business_name VARCHAR(200) NOT NULL,
    profession_type VARCHAR(50), -- 'auctioneer', 'appraiser', 'estate_sale_company'
    license_number VARCHAR(100),
    website_url TEXT,
    service_area TEXT, -- cities/regions served
    rating DECIMAL(3, 2) DEFAULT 0.00,
    review_count INTEGER DEFAULT 0,
    verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_professionals_user_id ON professionals(user_id);
CREATE INDEX idx_professionals_profession_type ON professionals(profession_type);
CREATE INDEX idx_professionals_verified ON professionals(verified);

-- Reviews (for professionals)
CREATE TABLE IF NOT EXISTS reviews (
    id SERIAL PRIMARY KEY,
    professional_id INTEGER REFERENCES professionals(id) ON DELETE CASCADE,
    reviewer_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    rating INTEGER CHECK (rating >= 1 AND rating <= 5),
    review_text TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_reviews_professional_id ON reviews(professional_id);
CREATE INDEX idx_reviews_reviewer_id ON reviews(reviewer_id);

-- Subscription plans (for sellers)
CREATE TABLE IF NOT EXISTS subscription_plans (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL, -- 'Basic', 'Pro', 'Premium'
    price DECIMAL(10, 2) NOT NULL,
    listings_per_month INTEGER,
    features JSONB, -- store feature flags as JSON
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- User subscriptions
CREATE TABLE IF NOT EXISTS user_subscriptions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    plan_id INTEGER REFERENCES subscription_plans(id),
    start_date TIMESTAMPTZ DEFAULT NOW(),
    end_date TIMESTAMPTZ,
    status VARCHAR(50) DEFAULT 'active', -- 'active', 'cancelled', 'expired'
    stripe_subscription_id TEXT, -- if using Stripe
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_user_subscriptions_user_id ON user_subscriptions(user_id);
CREATE INDEX idx_user_subscriptions_status ON user_subscriptions(status);

-- Insert default subscription plans
INSERT INTO subscription_plans (name, price, listings_per_month, features) VALUES
('Basic', 9.00, 3, '{"featured_listings": 0, "analytics": false, "priority_support": false}'),
('Pro', 19.00, 10, '{"featured_listings": 2, "analytics": true, "priority_support": false}'),
('Premium', 39.00, -1, '{"featured_listings": 5, "analytics": true, "priority_support": true}')
ON CONFLICT DO NOTHING;
