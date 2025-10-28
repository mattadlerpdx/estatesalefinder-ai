-- Seed script for sample estate sale data

-- First, ensure we have test users (sellers)
INSERT INTO users (firebase_uid, email, user_type, created_at)
VALUES
    ('test-seller-1', 'seller1@example.com', 'seller', NOW()),
    ('test-seller-2', 'seller2@example.com', 'seller', NOW()),
    ('test-buyer-1', 'buyer1@example.com', 'buyer', NOW())
ON CONFLICT (firebase_uid) DO NOTHING;

-- Get the user IDs for our sellers
DO $$
DECLARE
    seller1_id INT;
    seller2_id INT;
    sale1_id INT;
    sale2_id INT;
    sale3_id INT;
    sale4_id INT;
    sale5_id INT;
BEGIN
    -- Get seller IDs
    SELECT id INTO seller1_id FROM users WHERE firebase_uid = 'test-seller-1';
    SELECT id INTO seller2_id FROM users WHERE firebase_uid = 'test-seller-2';

    -- Insert estate sales
    INSERT INTO estate_sales (
        seller_id, title, description, sale_type, status,
        address_line1, address_line2, city, state, zip_code,
        latitude, longitude,
        start_date, end_date,
        listing_tier, featured, view_count,
        created_at, updated_at
    ) VALUES
    (
        seller1_id,
        'Vintage Furniture & Antiques Estate Sale',
        'Beautiful estate sale featuring vintage mid-century modern furniture, antique china, jewelry, books, and more! Everything must go. This charming home is filled with treasures from the 1950s-1970s. Don''t miss this opportunity to find unique pieces for your home.

Directions: From Highway 26, take exit 65. Turn left on NW 185th Ave, then right on NW Devoto Ln. House is on the left.

Parking: Street parking available. Please be courteous to neighbors.',
        'estate_sale',
        'published',
        '3773 NW Devoto Ln',
        NULL,
        'Portland',
        'OR',
        '97229',
        45.5428, -122.8089,
        NOW() + INTERVAL '2 days' + INTERVAL '9 hours',
        NOW() + INTERVAL '2 days' + INTERVAL '17 hours',
        'premium',
        true,
        127,
        NOW(),
        NOW()
    ) RETURNING id INTO sale1_id;

    INSERT INTO estate_sales (
        seller_id, title, description, sale_type, status,
        address_line1, city, state, zip_code,
        latitude, longitude,
        start_date, end_date,
        listing_tier, featured, view_count,
        created_at, updated_at
    ) VALUES
    (
        seller1_id,
        'Moving Sale - Tools, Electronics & Household Items',
        'We''re moving across the country and need to sell everything! Quality power tools, electronics, kitchen appliances, furniture, garden equipment, and more. All items priced to sell quickly. Cash only.',
        'moving_sale',
        'published',
        '1234 SE Hawthorne Blvd',
        'Portland',
        'OR',
        '97214',
        45.5122, -122.6519,
        NOW() + INTERVAL '5 days' + INTERVAL '10 hours',
        NOW() + INTERVAL '5 days' + INTERVAL '16 hours',
        'basic',
        false,
        89,
        NOW() - INTERVAL '2 days',
        NOW() - INTERVAL '2 days'
    ) RETURNING id INTO sale2_id;

    INSERT INTO estate_sales (
        seller_id, title, description, sale_type, status,
        address_line1, city, state, zip_code,
        latitude, longitude,
        start_date, end_date,
        listing_tier, featured, view_count,
        created_at, updated_at
    ) VALUES
    (
        seller2_id,
        'Downsizing Sale - Designer Clothes & Luxury Items',
        'High-end estate sale featuring designer clothing, handbags, shoes (sizes 6-8), fine jewelry, art pieces, and luxury home decor. Brands include Chanel, Louis Vuitton, Prada, and more. Serious buyers only.',
        'estate_sale',
        'published',
        '567 SW Broadway',
        'Portland',
        'OR',
        '97205',
        45.5202, -122.6796,
        NOW() + INTERVAL '1 day' + INTERVAL '11 hours',
        NOW() + INTERVAL '2 days' + INTERVAL '15 hours',
        'premium',
        true,
        234,
        NOW() - INTERVAL '5 days',
        NOW() - INTERVAL '5 days'
    ) RETURNING id INTO sale3_id;

    INSERT INTO estate_sales (
        seller_id, title, description, sale_type, status,
        address_line1, city, state, zip_code,
        latitude, longitude,
        start_date, end_date,
        listing_tier, featured, view_count,
        created_at, updated_at
    ) VALUES
    (
        seller2_id,
        'Complete Household Auction - Saturday 10am',
        'Live auction featuring complete contents of a beautiful Lake Oswego home. Furniture, appliances, artwork, collectibles, and more. Preview Friday 2-5pm. Auction starts Saturday at 10am sharp. Registration required.',
        'auction',
        'published',
        '890 Lake Shore Dr',
        'Lake Oswego',
        'OR',
        '97034',
        45.4207, -122.6709,
        NOW() + INTERVAL '3 days' + INTERVAL '10 hours',
        NOW() + INTERVAL '3 days' + INTERVAL '14 hours',
        'featured',
        false,
        156,
        NOW() - INTERVAL '1 day',
        NOW() - INTERVAL '1 day'
    ) RETURNING id INTO sale4_id;

    INSERT INTO estate_sales (
        seller_id, title, description, sale_type, status,
        address_line1, city, state, zip_code,
        latitude, longitude,
        start_date, end_date,
        listing_tier, featured, view_count,
        created_at, updated_at
    ) VALUES
    (
        seller1_id,
        'Multi-Family Garage Sale - Kids Toys & Baby Items',
        'Huge multi-family garage sale! Baby gear, kids toys, children''s clothing (newborn-10 years), strollers, car seats, books, games, and more. Great prices! Early birds welcome.',
        'garage_sale',
        'published',
        '2345 NE Glisan St',
        'Portland',
        'OR',
        '97232',
        45.5266, -122.6406,
        NOW() + INTERVAL '7 days' + INTERVAL '8 hours',
        NOW() + INTERVAL '7 days' + INTERVAL '13 hours',
        'basic',
        false,
        45,
        NOW(),
        NOW()
    ) RETURNING id INTO sale5_id;

    -- Insert sample images for the sales
    INSERT INTO sale_images (sale_id, image_url, is_primary, display_order)
    VALUES
        -- Sale 1 (Vintage Furniture) - Featured
        (sale1_id, 'https://images.unsplash.com/photo-1555041469-a586c61ea9bc?w=800', true, 1),
        (sale1_id, 'https://images.unsplash.com/photo-1493663284031-b7e3aefcae8e?w=800', false, 2),
        (sale1_id, 'https://images.unsplash.com/photo-1538688525198-9b88f6f53126?w=800', false, 3),

        -- Sale 2 (Moving Sale)
        (sale2_id, 'https://images.unsplash.com/photo-1484101403633-562f891dc89a?w=800', true, 1),
        (sale2_id, 'https://images.unsplash.com/photo-1581783898377-1c85bf937427?w=800', false, 2),

        -- Sale 3 (Designer Clothes) - Featured
        (sale3_id, 'https://images.unsplash.com/photo-1490481651871-ab68de25d43d?w=800', true, 1),
        (sale3_id, 'https://images.unsplash.com/photo-1567401893414-76b7b1e5a7a5?w=800', false, 2),
        (sale3_id, 'https://images.unsplash.com/photo-1591047139829-d91aecb6caea?w=800', false, 3),
        (sale3_id, 'https://images.unsplash.com/photo-1581044777550-4cfa60707c03?w=800', false, 4),

        -- Sale 4 (Auction)
        (sale4_id, 'https://images.unsplash.com/photo-1615529182904-14819c35db37?w=800', true, 1),
        (sale4_id, 'https://images.unsplash.com/photo-1616486338812-3dadae4b4ace?w=800', false, 2);

    -- Sale 5 (Garage Sale) has no images intentionally for variety

END $$;

-- Verify the data was inserted
SELECT
    id,
    LEFT(title, 40) as title,
    city,
    state,
    sale_type,
    featured,
    start_date::date as sale_date,
    view_count
FROM estate_sales
ORDER BY created_at DESC;

SELECT COUNT(*) as total_sales FROM estate_sales;
SELECT COUNT(*) as total_images FROM sale_images;
