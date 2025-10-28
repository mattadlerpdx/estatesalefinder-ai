# Competitor Flow Analysis - EstateSale-Finder.com

## Screenshots Location
`docs/research/competitor-screenshots/`

## User Flow Analysis

### 1. Dashboard Home (Screenshot 01)
**Page**: Home Page for Cami (User Dashboard)

**Features:**
- **"Add Sale" button** - Primary CTA at top
- **Balance display** - Shows user owes $50 (pay-as-you-go model)
- **"Make a PayPal Payment" button** - Payment integration
- **Recent Transactions** - Shows charge of $50 for Estate Sale ID# 15437
- **Three sections:**
  - Current Sales (0 sales)
  - Upcoming Sales (0 sales)
  - Past Sales (1 sale shown)
- **Past sale listing shows:**
  - Sale ID# 15437
  - Address: 3773 nw devoto ln, Yes OR 97229
  - Date: 26th Oct 9:00am
  - Type: 1 day Estate Sale
  - Signup count
  - "View" button

**Our Improvements:**
- Modern card-based UI instead of tables
- Real-time balance (no owing money model)
- Clearer sale status indicators
- Quick actions (edit, delete, duplicate)

---

### 2. Add A Sale Form (Screenshot 04)
**Page**: Add A Sale

**Form Sections:**

#### Sale Location
- Address line 1 (required)
- Address line 2 (optional)
- City/State/Zip dropdown (Oregon selected)
- Region selector
- Special Driving/Parking Directions (rich text editor)

#### Sale Details
- **Choose Sale Frequency:**
  - One Time (radio button)
  - Recurring (radio button)
- **Date of Sale:** 10/26/2025 (date picker)
- **Start Time:** 09:00 AM (time picker)
- **End Date:** 10/26/2025 (for multi-day)
- **End Time:** 05:00 PM (time picker)
- **Sale Type:** Estate Sale (dropdown)
- **Sale Text:** Rich text editor (for description)
  - Note: "Remember to add any sale specific contact information here"

#### Signup Details
- Type: None (dropdown)

#### OPTIONAL: Featured Sale of the Week
- Checkbox: "I would like to list this sale on the ESF homepage as a Featured Sale of the Week. By checking this box, I agree to pay the Featured Sale listing fee of $100. Featured Sales of the week are visible on the front page of ESF from 12:01 am on Monday of the sale week until the following Sunday at 11:59 pm or the conclusion of the sale, whichever is first."

#### Terms Agreement
- Checkbox: "I agree to the ESF terms and conditions, and acknowledge that I am responsible for the listing fee."
- Note: "ESF is not responsible for ad content and/or errors. Please verify your sale information before clicking 'next'."

**"Next" button** at bottom

**Our Improvements:**
- Multi-step wizard instead of long form
- Auto-save drafts
- Address autocomplete with Google Maps
- Image upload in same flow (not separate)
- Real-time preview
- Better mobile experience

---

### 3. Sale Display Time Settings (Screenshot 03)
**Page**: Sale Submitted

**Content:**
- **Success message:** "You sale has been submitted and will be activated shortly. The following screens will enable you to control when the sale is displayed and also to add photographs to your sale listing."

#### Sale Display Time
- **Explanation:** "You can specify when you want your sale details listed. The default is to show the details 24hrs before the sale starts and stop displaying as soon as the sale ends. You can change these values below."

#### Posting Details
- **When do you want the address of the sale posted?**
  - Start Date: 10/25/2025 (date picker)
  - Start Time: 09:00 AM (time picker)
  - Note: "If no time selector shown use 24hr format e.g 16:00 for 4pm"

**"Next" button**

**Our Improvements:**
- Make this optional (smart defaults)
- Combine with main form
- Clearer UX - most users won't need this
- Auto-calculate posting dates

---

### 4. Add Photos (Screenshot 02)
**Page**: Add Photos for Sale Number 15437

**Features:**
- **"Choose Files" button** - File upload
- **Instructions:**
  - Select multiple photos with shift or control keys
  - If experiencing issues:
    - Try uploading photos one at a time
    - Be patient - uploading multiple photos can take a while
    - Try a different browser (e.g., Chrome instead of IE)
    - Scale photos before uploading (modern cameras = high res = very large files)

**Action Buttons:**
- "Upload Photos"
- "Done for Now"

**Our Improvements:**
- Drag & drop interface
- Image preview before upload
- Automatic image compression
- Progress indicators
- Set primary image
- Reorder images
- Image editing tools (crop, rotate)
- Much better UX!

---

## Key Insights

### Their Model:
1. **Pay-as-you-go** - User owes $50 after posting
2. **Separate steps** - Sale creation → Display time → Photo upload
3. **Manual activation** - "will be activated shortly" (admin approval?)
4. **Featured listings** - $100 upgrade for homepage placement
5. **PayPal integration** - For payments

### Our Competitive Advantages:
1. **Cleaner UX** - Modern Tailwind design
2. **Single flow** - Create sale with photos in one go
3. **Instant publishing** - No manual approval needed
4. **Cheaper pricing** - $9/mo vs $50 per sale
5. **Better mobile** - Responsive Next.js app
6. **Auto-save** - Never lose progress
7. **AI features** - Smart categorization, price suggestions

---

## Recommended Our Flow

### Step 1: Sale Details
- Title
- Description (rich text)
- Sale type dropdown
- Date range picker
- Hours/schedule

### Step 2: Location
- Address autocomplete
- Map preview
- Parking/directions (optional)

### Step 3: Photos
- Drag & drop upload
- Set primary image
- Reorder images
- Auto-compress

### Step 4: Review & Publish
- Preview card
- Edit any section
- Publish immediately or schedule
- Choose listing tier (basic/featured/premium)

### Result: Dashboard
- See your sale live
- Edit anytime
- View analytics
- Manage photos

---

## Next Steps for Development

1. **Build `/sales` page first** - Public browsing
2. **Build `/dashboard` page** - Seller management
3. **Build create sale flow** - Multi-step form
4. **Add image upload** - Drag & drop UI
5. **Add payment** - Stripe integration

---

**Key Takeaway:** Their UX is dated and clunky. We can do WAY better with modern Next.js + Tailwind CSS!
