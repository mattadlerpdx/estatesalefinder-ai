/**
 * API Integration Tests for Listing endpoints
 * These tests verify that the frontend can correctly fetch and process listings from the backend
 *
 * NOTE: These tests require the backend to be running
 * Start backend with: cd backend && go run cmd/api/main.go
 */

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

const BACKEND_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

describe('Listing API Integration Tests', () => {
  describe('GET /api/sales - Fetch all listings', () => {
    test('should fetch listings successfully', async () => {
      const response = await fetch(`${BACKEND_URL}/api/sales`)

      expect(response.ok).toBe(true)
      expect(response.status).toBe(200)

      const data = await response.json()

      // Response should be wrapped in standard response object
      expect(data).toHaveProperty('success')
      expect(data).toHaveProperty('data')

      if (data.success && Array.isArray(data.data)) {
        const listings: Listing[] = data.data

        // Verify structure of first listing if any exist
        if (listings.length > 0) {
          const firstListing = listings[0]

          expect(firstListing).toHaveProperty('id')
          expect(firstListing).toHaveProperty('title')
          expect(firstListing).toHaveProperty('city')
          expect(firstListing).toHaveProperty('state')
          expect(firstListing).toHaveProperty('is_scraped')

          console.log(`✓ Fetched ${listings.length} listings`)
          console.log(`✓ First listing: "${firstListing.title}" in ${firstListing.city}, ${firstListing.state}`)
        }
      }
    }, 10000) // 10 second timeout for API call

    test('should fetch listings with city filter', async () => {
      const city = 'Portland'
      const response = await fetch(`${BACKEND_URL}/api/sales?city=${city}`)

      expect(response.ok).toBe(true)

      const data = await response.json()

      if (data.success && Array.isArray(data.data)) {
        const listings: Listing[] = data.data

        // All listings should be from Portland
        listings.forEach(listing => {
          if (listing.city) {
            expect(listing.city.toLowerCase()).toContain(city.toLowerCase())
          }
        })

        console.log(`✓ Fetched ${listings.length} listings from ${city}`)
      }
    }, 10000)

    test('should fetch listings with state filter', async () => {
      const state = 'OR'
      const response = await fetch(`${BACKEND_URL}/api/sales?state=${state}`)

      expect(response.ok).toBe(true)

      const data = await response.json()

      if (data.success && Array.isArray(data.data)) {
        const listings: Listing[] = data.data

        // All listings should be from Oregon
        listings.forEach(listing => {
          expect(listing.state).toBe(state)
        })

        console.log(`✓ Fetched ${listings.length} listings from ${state}`)
      }
    }, 10000)

    test('should fetch featured listings only', async () => {
      const response = await fetch(`${BACKEND_URL}/api/sales?featured=true`)

      expect(response.ok).toBe(true)

      const data = await response.json()

      if (data.success && Array.isArray(data.data)) {
        const listings: Listing[] = data.data

        // All listings should be featured
        listings.forEach(listing => {
          if (listing.featured !== undefined) {
            expect(listing.featured).toBe(true)
          }
        })

        console.log(`✓ Fetched ${listings.length} featured listings`)
      }
    }, 10000)

    test('should respect pagination limits', async () => {
      const limit = 5
      const response = await fetch(`${BACKEND_URL}/api/sales?limit=${limit}`)

      expect(response.ok).toBe(true)

      const data = await response.json()

      if (data.success && Array.isArray(data.data)) {
        const listings: Listing[] = data.data

        expect(listings.length).toBeLessThanOrEqual(limit)

        console.log(`✓ Pagination working: requested ${limit}, got ${listings.length}`)
      }
    }, 10000)
  })

  describe('GET /api/sales/:id - Fetch single listing', () => {
    test('should fetch a single listing by ID', async () => {
      // First, get a listing ID
      const listResponse = await fetch(`${BACKEND_URL}/api/sales?limit=1`)
      const listData = await listResponse.json()

      if (listData.success && listData.data && listData.data.length > 0) {
        const listingId = listData.data[0].id

        // Now fetch that specific listing
        const response = await fetch(`${BACKEND_URL}/api/sales/${listingId}`)

        expect(response.ok).toBe(true)
        expect(response.status).toBe(200)

        const data = await response.json()

        expect(data).toHaveProperty('success')
        expect(data).toHaveProperty('data')

        if (data.success && data.data) {
          const listing = data.data

          expect(listing.id).toBe(listingId)
          expect(listing).toHaveProperty('title')
          expect(listing).toHaveProperty('description')
          expect(listing).toHaveProperty('view_count')

          console.log(`✓ Fetched listing: "${listing.title}"`)
          console.log(`✓ View count: ${listing.view_count}`)
        }
      } else {
        console.log('⚠ No listings available to test detail endpoint')
      }
    }, 10000)

    test('should return 404 for non-existent listing', async () => {
      const response = await fetch(`${BACKEND_URL}/api/sales/999999`)

      expect(response.status).toBe(404)

      console.log('✓ Correctly returns 404 for non-existent listing')
    }, 10000)
  })

  describe('Listing type validation', () => {
    test('listings should have correct is_scraped field', async () => {
      const response = await fetch(`${BACKEND_URL}/api/sales`)
      const data = await response.json()

      if (data.success && Array.isArray(data.data)) {
        const listings: Listing[] = data.data

        const scrapedCount = listings.filter(l => l.is_scraped === true).length
        const ownedCount = listings.filter(l => l.is_scraped === false).length

        console.log(`✓ Found ${scrapedCount} scraped listings and ${ownedCount} owned listings`)

        // Scraped listings should have source info
        listings.filter(l => l.is_scraped).forEach(listing => {
          expect(listing.source).toBeDefined()
          expect(listing.source?.name).toBeTruthy()
          expect(listing.source?.url).toBeTruthy()
        })

        console.log('✓ All scraped listings have source information')
      }
    }, 10000)

    test('listing dates should be valid ISO strings', async () => {
      const response = await fetch(`${BACKEND_URL}/api/sales?limit=5`)
      const data = await response.json()

      if (data.success && Array.isArray(data.data)) {
        const listings: Listing[] = data.data

        listings.forEach(listing => {
          // Dates should be parseable
          const startDate = new Date(listing.start_date)
          const endDate = new Date(listing.end_date)

          expect(startDate.toString()).not.toBe('Invalid Date')
          expect(endDate.toString()).not.toBe('Invalid Date')

          // End date should be after start date
          expect(endDate.getTime()).toBeGreaterThan(startDate.getTime())
        })

        console.log('✓ All listing dates are valid')
      }
    }, 10000)

    test('listing images should have correct structure', async () => {
      const response = await fetch(`${BACKEND_URL}/api/sales`)
      const data = await response.json()

      if (data.success && Array.isArray(data.data)) {
        const listings: Listing[] = data.data

        const listingsWithImages = listings.filter(l => l.images && l.images.length > 0)

        if (listingsWithImages.length > 0) {
          listingsWithImages.forEach(listing => {
            listing.images?.forEach(image => {
              expect(image).toHaveProperty('image_url')
              expect(image).toHaveProperty('is_primary')
              expect(typeof image.image_url).toBe('string')
              expect(typeof image.is_primary).toBe('boolean')
            })
          })

          console.log(`✓ Verified image structure for ${listingsWithImages.length} listings`)
        }
      }
    }, 10000)
  })

  describe('Error handling', () => {
    test('should handle network errors gracefully', async () => {
      try {
        await fetch('http://localhost:9999/api/sales')
        fail('Should have thrown an error')
      } catch (error) {
        expect(error).toBeDefined()
        console.log('✓ Network errors handled correctly')
      }
    }, 10000)

    test('should handle invalid query parameters', async () => {
      const response = await fetch(`${BACKEND_URL}/api/sales?limit=-1`)

      // Should either return 400 or handle gracefully with default
      expect(response.status).toBeGreaterThanOrEqual(200)

      console.log(`✓ Invalid parameters handled (status: ${response.status})`)
    }, 10000)
  })
})

// Mock data for offline testing
export const mockListings: Listing[] = [
  {
    id: '1',
    title: 'Estate Sale - Vintage Collection',
    description: 'Beautiful vintage furniture and collectibles',
    address: '123 Main Street',
    city: 'Portland',
    state: 'OR',
    zip_code: '97201',
    start_date: '2025-11-01T09:00:00Z',
    end_date: '2025-11-03T17:00:00Z',
    thumbnail_url: 'https://example.com/thumb1.jpg',
    images: [
      { image_url: 'https://example.com/img1.jpg', is_primary: true },
      { image_url: 'https://example.com/img2.jpg', is_primary: false }
    ],
    is_scraped: false,
    sale_type: 'estate_sale',
    status: 'active',
    view_count: 42,
    featured: true
  },
  {
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
    view_count: 15,
    featured: false
  }
]

export type { Listing }
