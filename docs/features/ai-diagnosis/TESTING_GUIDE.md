# Testing AI Attachment Analysis

## Quick Setup

### 1. Add OpenAI API Key
Edit .env file:
\\\ash
OPENAI_API_KEY=sk-proj-your-actual-openai-api-key-here
\\\

### 2. Restart Backend
\\\powershell
cd C:\Users\birju\aby-med
.\bin\platform.exe
\\\

### 3. Restart Frontend
\\\powershell
cd C:\Users\birju\aby-med\admin-ui
npm run dev
\\\

## Testing on Service Request Page

1. **Navigate to**: http://localhost:3002/service-request?qr=QR-MAP-0001
2. **Fill in**:
   - Your Name: Test User
   - Priority: Medium
   - Description: Equipment display showing error code E-042
3. **Attach Files**: Upload 1-2 images or a video
4. **Click**: "Get AI Diagnosis with 2 Image(s)" button
5. **Watch Progress**:
   - Processing image 1/2...
   - Analyzing 2 image(s) with AI...
   - Analysis complete!
6. **View Results**: Check console for diagnosis output

## What the AI Analyzes

### From Images:
- Error codes visible on displays
- Physical damage (cracks, burns, corrosion)
- Component identification
- Wear patterns

### From Videos:
- Extracts 3 frames automatically
- Analyzes each frame
- Detects unusual behavior

## Expected Behavior

### Success Scenario:
- Button shows progress messages
- Green success indicator after complete
- Diagnosis object in console with:
  - primary_diagnosis
  - vision_analysis
  - recommended_actions
  - required_parts

### Error Scenarios:
- No OpenAI key → "API key not configured"
- No description → "Please provide a description first"
- API error → Red error message displayed

## Features Implemented

✅ Image compression before sending
✅ Video frame extraction (3 frames)
✅ Progress indicators
✅ Optional/independent from ticket creation
✅ Error handling
✅ Visual feedback

## Notes

- AI analysis is **completely optional**
- Works independently - doesn't block ticket creation
- Can test without submitting ticket
- Diagnosis results logged to console
- Visual analysis included if files uploaded

---
*Created: December 11, 2025*
