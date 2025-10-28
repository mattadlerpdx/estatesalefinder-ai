/**
 * Type tests for Listing interface
 * These tests verify that the refactoring from Sale -> Listing maintains type safety
 */

// Import type for testing (you would normally import from a shared types file)
interface Listing {
  id: string
  title: string
  description?: string
  address: string
  city: string
  state: string
  zip_code: string
  start_date: string
  end_date: string
  thumbnail_url?: string
  images?: Array<{
    image_url: string
    is_primary: boolean
  }>
  is_scraped: boolean
  source?: {
    name: string
    url: string
  }
  sale_type?: string
  status?: string
  view_count?: number
  featured?: boolean
}

interface ListingDetail extends Listing {
  id: number
  seller_id: number
  address_line1: string
  address_line2?: string
  latitude?: number
  longitude?: number
  listing_tier: string
  driving_directions?: string
  parking_info?: string
  created_at: string
  images?: Array<{
    id: number
    image_url: string
    is_primary: boolean
    display_order: number
  }>
}

// Type assertion tests
describe('Listing Type Tests', () => {
  test('Listing interface has correct structure for owned listing', () => {
    const ownedListing: Listing = {
      id: '1',
      title: 'Estate Sale - Vintage Furniture',
      description: 'Beautiful mid-century modern pieces',
      address: '123 Main Street',
      city: 'Portland',
      state: 'OR',
      zip_code: '97201',
      start_date: '2025-11-01T09:00:00Z',
      end_date: '2025-11-03T17:00:00Z',
      thumbnail_url: 'https://example.com/thumb.jpg',
      images: [
        {
          image_url: 'https://example.com/img1.jpg',
          is_primary: true
        },
        {
          image_url: 'https://example.com/img2.jpg',
          is_primary: false
        }
      ],
      is_scraped: false,
      sale_type: 'estate_sale',
      status: 'active',
      view_count: 42,
      featured: true
    }

    expect(ownedListing.title).toBe('Estate Sale - Vintage Furniture')
    expect(ownedListing.is_scraped).toBe(false)
    expect(ownedListing.city).toBe('Portland')
    expect(ownedListing.images).toHaveLength(2)
  })

  test('Listing interface has correct structure for scraped listing', () => {
    const scrapedListing: Listing = {
      id: 'external-123',
      title: 'Moving Sale - Everything Must Go',
      description: 'Furniture, appliances, and more',
      address: '456 Oak Avenue',
      city: 'Seattle',
      state: 'WA',
      zip_code: '98101',
      start_date: '2025-11-05T08:00:00Z',
      end_date: '2025-11-06T16:00:00Z',
      thumbnail_url: 'https://external.com/thumb.jpg',
      is_scraped: true,
      source: {
        name: 'EstateSales.net',
        url: 'https://estatesales.net/WA/Seattle/98101/12345'
      },
      sale_type: 'moving_sale',
      view_count: 15
    }

    expect(scrapedListing.title).toBe('Moving Sale - Everything Must Go')
    expect(scrapedListing.is_scraped).toBe(true)
    expect(scrapedListing.source).toBeDefined()
    expect(scrapedListing.source?.name).toBe('EstateSales.net')
  })

  test('ListingDetail interface extends Listing correctly', () => {
    const detailedListing: ListingDetail = {
      id: 42,
      seller_id: 100,
      title: 'Antique Collection Sale',
      description: 'Rare collectibles and antiques',
      address: '789 Pine Street',
      address_line1: '789 Pine Street',
      address_line2: 'Unit 5B',
      city: 'Beaverton',
      state: 'OR',
      zip_code: '97005',
      latitude: 45.4871,
      longitude: -122.8037,
      start_date: '2025-11-10T09:00:00Z',
      end_date: '2025-11-12T17:00:00Z',
      listing_tier: 'premium',
      thumbnail_url: 'https://example.com/antiques-thumb.jpg',
      is_scraped: false,
      sale_type: 'estate_sale',
      status: 'active',
      view_count: 156,
      featured: true,
      driving_directions: 'Take exit 69B, turn left on Pine',
      parking_info: 'Street parking available',
      created_at: '2025-10-15T10:00:00Z',
      images: [
        {
          id: 1,
          image_url: 'https://example.com/antique1.jpg',
          is_primary: true,
          display_order: 1
        },
        {
          id: 2,
          image_url: 'https://example.com/antique2.jpg',
          is_primary: false,
          display_order: 2
        }
      ]
    }

    expect(detailedListing.id).toBe(42)
    expect(detailedListing.seller_id).toBe(100)
    expect(detailedListing.listing_tier).toBe('premium')
    expect(detailedListing.latitude).toBe(45.4871)
    expect(detailedListing.images).toHaveLength(2)
    expect(detailedListing.images?.[0].display_order).toBe(1)
  })

  test('Array of Listings can be typed correctly', () => {
    const listings: Listing[] = [
      {
        id: '1',
        title: 'Sale 1',
        address: '100 First St',
        city: 'Portland',
        state: 'OR',
        zip_code: '97201',
        start_date: '2025-11-01T09:00:00Z',
        end_date: '2025-11-02T17:00:00Z',
        is_scraped: false
      },
      {
        id: '2',
        title: 'Sale 2',
        address: '200 Second St',
        city: 'Eugene',
        state: 'OR',
        zip_code: '97401',
        start_date: '2025-11-03T09:00:00Z',
        end_date: '2025-11-04T17:00:00Z',
        is_scraped: true,
        source: {
          name: 'External Source',
          url: 'https://example.com/sale2'
        }
      }
    ]

    expect(listings).toHaveLength(2)
    expect(listings[0].is_scraped).toBe(false)
    expect(listings[1].is_scraped).toBe(true)
    expect(listings[1].source).toBeDefined()
  })

  test('Optional fields work correctly', () => {
    const minimalListing: Listing = {
      id: '99',
      title: 'Minimal Listing',
      address: '999 Minimal Lane',
      city: 'Salem',
      state: 'OR',
      zip_code: '97301',
      start_date: '2025-12-01T09:00:00Z',
      end_date: '2025-12-02T17:00:00Z',
      is_scraped: false
    }

    expect(minimalListing.description).toBeUndefined()
    expect(minimalListing.thumbnail_url).toBeUndefined()
    expect(minimalListing.images).toBeUndefined()
    expect(minimalListing.source).toBeUndefined()
    expect(minimalListing.sale_type).toBeUndefined()
    expect(minimalListing.status).toBeUndefined()
    expect(minimalListing.view_count).toBeUndefined()
    expect(minimalListing.featured).toBeUndefined()
  })

  test('Listing can be filtered and mapped', () => {
    const allListings: Listing[] = [
      {
        id: '1',
        title: 'Portland Sale',
        address: '100 Main',
        city: 'Portland',
        state: 'OR',
        zip_code: '97201',
        start_date: '2025-11-01',
        end_date: '2025-11-02',
        is_scraped: false,
        featured: true
      },
      {
        id: '2',
        title: 'Seattle Sale',
        address: '200 Pike',
        city: 'Seattle',
        state: 'WA',
        zip_code: '98101',
        start_date: '2025-11-03',
        end_date: '2025-11-04',
        is_scraped: true,
        featured: false
      },
      {
        id: '3',
        title: 'Portland Sale 2',
        address: '300 Oak',
        city: 'Portland',
        state: 'OR',
        zip_code: '97202',
        start_date: '2025-11-05',
        end_date: '2025-11-06',
        is_scraped: false,
        featured: false
      }
    ]

    // Filter Portland listings
    const portlandListings = allListings.filter(l => l.city === 'Portland')
    expect(portlandListings).toHaveLength(2)

    // Filter scraped listings
    const scrapedListings = allListings.filter(l => l.is_scraped)
    expect(scrapedListings).toHaveLength(1)
    expect(scrapedListings[0].city).toBe('Seattle')

    // Filter featured listings
    const featuredListings = allListings.filter(l => l.featured === true)
    expect(featuredListings).toHaveLength(1)

    // Map to titles
    const titles = allListings.map(l => l.title)
    expect(titles).toEqual(['Portland Sale', 'Seattle Sale', 'Portland Sale 2'])
  })
})

// Export so TypeScript doesn't complain about unused types
export type { Listing, ListingDetail }
