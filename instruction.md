I want to make a rentry.co clone called Pretrey Next



The tech stack is

Go

HTML

CSS

Lighting Design system 2

SQLite



For logins I dont want to use like email sign ons it is simple login with username and pass and we use hash cash to prevent spam.

For signups the same



And also use hashcash when creating the page

and also have a admin dashboard on /admin and its protected by username and password in the env (this is not shown to the public this page) the admin dash support deleting pages search pages show statistics etc.

the project is called Zwidth Paste and dont call it rentry clone

btw here is the spec sheet to follow

# Rentry.co - Feature Specification Document

**Version:** 1.0  
**Date:** February 17, 2026  
**Purpose:** User flows and feature overview for recreating Rentry.co

---

## Table of Contents

1. [Product Overview](#product-overview)
2. [Core User Flows](#core-user-flows)
3. [Feature List](#feature-list)
4. [Page Specifications](#page-specifications)
5. [Content Formatting](#content-formatting)
6. [Access & Permissions](#access--permissions)

---

## Product Overview

### What is Rentry?

Rentry is a simple, fast markdown-based text sharing platform where users can:

- Create and share text content instantly
- Use custom or random URLs
- Protect content with edit codes
- Preview content in real-time
- Export to various formats

### Key Principles

| Principle    | Description                                     |
| ------------ | ----------------------------------------------- |
| **Simple**   | No registration required for basic use          |
| **Fast**     | Instant creation and editing                    |
| **Free**     | Core features available at no cost              |
| **Flexible** | Custom URLs, edit codes, and formatting options |

---

## Core User Flows

### Flow 1: Create New Entry

```
User visits homepage
    ↓
Enters markdown content in editor
    ↓
(Optional) Sets custom URL
    ↓
(Optional) Sets edit code
    ↓
Clicks "Go" or presses Ctrl+Enter
    ↓
Entry is created
    ↓
User sees:
  - Entry URL
  - Edit code (if generated)
  - Option to edit immediately
```

**Notes:**

- If no URL provided, system generates random URL
- If no edit code provided, system generates random code
- Edit code shown only once - user must save it

---

### Flow 2: View Entry

```
User visits entry URL (rentry.co/{url})
    ↓
System loads entry content
    ↓
Content rendered as formatted HTML
    ↓
User sees:
  - Formatted content
  - Publication date
  - Last edit date
  - View count
  - Edit button (if has edit code)
  - Export button
```

---

### Flow 3: Edit Entry

```
User clicks "Edit" button or visits /edit/{url}
    ↓
Enters edit code
    ↓
System verifies code
    ↓
Editor loads with existing content
    ↓
User makes changes
    ↓
(Optional) Changes URL
    ↓
(Optional) Changes edit code
    ↓
(Optional) Sets modify code
    ↓
Clicks "Save"
    ↓
Entry updated, redirects to entry page
```

---

### Flow 4: Delete Entry

```
User visits /edit/{url}
    ↓
Enters edit code
    ↓
Clicks "Delete" button
    ↓
Confirms deletion
    ↓
Entry marked as deleted
    ↓
URL becomes available for reuse
```

---

### Flow 5: User Login (Optional)

```
User clicks "Login"
    ↓
Enters email address
    ↓
System sends login link to email
    ↓
User clicks link in email
    ↓
User logged in (session created)
    ↓
Enhanced features unlocked:
  - Captcha-free exports
  - Claim credits
  - API access
```

---

## Feature List

### P0 - Must Have (MVP)

| Feature             | Description                                      | Priority |
| ------------------- | ------------------------------------------------ | -------- |
| **Markdown Editor** | Text area for entering markdown content          | P0       |
| **Preview Mode**    | Toggle between raw markdown and rendered view    | P0       |
| **Random URLs**     | Auto-generate unique URL slugs                   | P0       |
| **Custom URLs**     | Allow users to choose their own URL              | P0       |
| **Edit Codes**      | Password-protected editing                       | P0       |
| **Basic Markdown**  | Headers, bold, italic, links, lists, code blocks | P0       |
| **Entry Deletion**  | Delete entries with edit code                    | P0       |

### P1 - Should Have

| Feature                 | Description                                    | Priority |
| ----------------------- | ---------------------------------------------- | -------- |
| **Modify Codes**        | Limited access codes (text editing only)       | P1       |
| **Syntax Highlighting** | Code block highlighting for multiple languages | P1       |
| **Image Support**       | Inline images with sizing options              | P1       |
| **Table of Contents**   | Auto-generated from headers                    | P1       |
| **View Counter**        | Track and display page views                   | P1       |
| **Export (PDF)**        | Export entries as PDF documents                | P1       |
| **Metadata**            | Custom key-value storage for entries           | P1       |
| **Passwordless Login**  | Email-based authentication                     | P1       |

### P2 - Nice to Have

| Feature              | Description                       | Priority |
| -------------------- | --------------------------------- | -------- |
| **Colored Text**     | Custom text colors                | P2       |
| **Underlined Text**  | Custom underline styles           | P2       |
| **Admonitions**      | Info/warning/danger callout boxes | P2       |
| **Spoilers**         | Hidden text that reveals on click | P2       |
| **Export (Images)**  | PNG, JPG, WebP export options     | P2       |
| **URL Claim System** | Claim inactive URLs               | P2       |
| **API Access**       | REST API for programmatic access  | P2       |
| **CLI Tool**         | Command-line interface            | P2       |
| **Rate Limiting**    | Prevent abuse                     | P2       |

---

## Page Specifications

### Homepage (/)

**Purpose:** Create new entries

**Components:**

- Tab navigation (Text / Preview / How)
- Large text editor area
- Metadata toggle button
- Custom edit code input
- Custom URL input
- "Go" button
- Navigation links (new, what, how, langs, contacts, login)

**Actions:**

- Create entry with content
- Set optional custom URL
- Set optional edit code
- Set optional metadata

---

### Entry Page (/{url})

**Purpose:** View rendered entry content

**Components:**

- Edit button (top right)
- Export button (top right)
- Rendered markdown content
- Publication timestamp
- Last edit timestamp
- View count
- Navigation links (footer)

**Actions:**

- View content
- Click edit (requires edit code)
- Export entry
- Click header permalinks

---

### Edit Page (/edit/{url})

**Purpose:** Modify existing entries

**Components:**

- Tab navigation (Text / Preview / How)
- Text editor (pre-filled)
- Edit code input field
- New edit code input (optional)
- New URL input (optional)
- New modify code input (optional)
- Save, Back, Delete buttons

**Actions:**

- Update content
- Change edit code
- Change URL
- Set modify code
- Delete entry

---

### About Page (/what)

**Purpose:** Explain service features

**Content:**

- Service description
- URL system explanation
- Edit code explanation
- Modify code explanation
- Character limits
- Rules and guidelines
- Duration policy
- Contact information

---

### Cheatsheet Page (/how)

**Purpose:** Markdown syntax reference

**Content:**

- Header examples
- Text formatting (bold, italic, strike, mark)
- Colored text examples
- Underlined text examples
- Link and image examples
- List examples (ordered, unordered, checkbox)
- Code block examples
- Table examples
- Blockquote examples
- Special features (TOC, admonitions, spoilers)

---

### Languages Page (/langs)

**Purpose:** List supported syntax highlighting languages

**Content:**

- Code block example
- Complete list of 200+ language names and aliases

---

### Login Page (/request-login)

**Purpose:** User authentication

**Components:**

- Email input field
- "Remember me" checkbox
- "Send Login Link" button

**Flow:**

1. User enters email
2. System sends magic link
3. User clicks link
4. Session created

---

## Content Formatting

### Basic Markdown

| Feature       | Syntax          | Example               |
| ------------- | --------------- | --------------------- |
| Headers       | `#` to `######` | `# Heading`           |
| Bold          | `**text**`      | `**bold**`            |
| Italic        | `*text*`        | `*italic*`            |
| Strikethrough | `~~text~~`      | `~~deleted~~`         |
| Highlight     | `==text==`      | `==marked==`          |
| Links         | `[text](url)`   | `[link](https://...)` |
| Images        | `![alt](url)`   | `![img](url)`         |
| Code (inline) | `` `code` ``    | `` `inline` ``        |
| Code (block)  | ` ```lang `     | ` ```python `         |
| Lists         | `- item`        | `- bullet`            |
| Checkboxes    | `- [ ]`         | `- [ ] task`          |
| Quotes        | `> text`        | `> quoted`            |
| Tables        | `\| col \|`     | `\| a \| b \|`        |

### Extended Formatting

| Feature           | Syntax          | Description                    |
| ----------------- | --------------- | ------------------------------ |
| Colored Text      | `%red% text %%` | Named or hex colors            |
| Underline         | `!~ text ~!`    | With optional style/color      |
| Centered          | `-> text <-`    | Center alignment               |
| Right Align       | `-> text`       | Right alignment                |
| Table of Contents | `[TOC]`         | Auto-generated from headers    |
| Admonitions       | `!!! type`      | Info/note/warning/danger boxes |
| Spoilers          | `!> text`       | Hidden until clicked           |
| Comments          | `[//]: (text)`  | Not rendered                   |

---

## Access & Permissions

### Access Levels

| Level                    | Capabilities                          | How to Obtain     |
| ------------------------ | ------------------------------------- | ----------------- |
| **Viewer**               | View content, export                  | Anyone with URL   |
| **Editor (Modify Code)** | Edit text only                        | Given modify code |
| **Owner (Edit Code)**    | Full control (edit, delete, transfer) | Given edit code   |

### Edit Code vs Modify Code

| Action             | Edit Code | Modify Code |
| ------------------ | --------- | ----------- |
| Edit text          | ✅         | ✅           |
| Change URL         | ✅         | ❌           |
| Change edit code   | ✅         | ❌           |
| Change modify code | ✅         | ❌           |
| Delete entry       | ✅         | ❌           |
| View metadata      | ✅         | ❌           |

### Recovery Options

If edit code is lost:

1. **With verification setup:**
   - Email support from verified email
   - Prove ownership via linked social account

2. **Without verification:**
   - Contact support with evidence
   - Provide creation IP, date, content details
   - No guarantee of recovery

**Best Practice:** Set `SECRET_EMAIL_ADDRESS` in metadata when creating important entries.

---

## Limits & Constraints

### Content Limits

| Resource     | Limit              |
| ------------ | ------------------ |
| Text content | 200,000 characters |
| Custom URL   | 2-100 characters   |
| Edit code    | 1-100 characters   |
| Modify code  | 1-100 characters   |

### URL Rules

- Allowed characters: `a-z`, `0-9`, `_` (underscore), `-` (hyphen)
- Minimum length: 2 characters
- Maximum length: 100 characters
- Case: Stored lowercase, original case preserved

### Rate Limits

| Action             | Limit                |
| ------------------ | -------------------- |
| URL claims         | 3 per week per email |
| Export (anonymous) | ~5 per hour          |
| Export (logged in) | ~50 per hour         |
| Login requests     | ~5 per hour          |

---

## Export Options

### Available Formats

| Format   | Use Case            | Notes            |
| -------- | ------------------- | ---------------- |
| **PDF**  | Documents, printing | Best for text    |
| **PNG**  | Images, screenshots | Lossless         |
| **JPG**  | Web images          | Compressed       |
| **WebP** | Modern web          | Best compression |

### Export Settings

- Quality (for image formats)
- Page width
- Theme (light/dark/auto)

---

## Metadata Fields

### Standard Fields

| Field                      | Purpose                   | Visible |
| -------------------------- | ------------------------- | ------- |
| `SECRET_EMAIL_ADDRESS`     | Ownership verification    | No      |
| `SECRET_VERIFY`            | Social media verification | No      |
| `SECRET_RAW_ACCESS_CODE`   | API access without auth   | No      |
| `ACCESS_RECOMMENDED_THEME` | Suggested display theme   | Yes     |

### Metadata Rules

- SECRET fields cannot be viewed after setting
- SECRET fields require edit code (not modify code)
- Maximum 50 metadata entries per entry
- Keys: max 100 characters
- Values: max 1000 characters

---

## URL Claim System

### When URLs Can Be Claimed

| Scenario            | Can Claim? | Waiting Period                |
| ------------------- | ---------- | ----------------------------- |
| Never edited        | Yes        | 12 months inactive            |
| Edited recently     | No         | -                             |
| Previously claimed  | Yes        | 3 months + no genuine content |
| Has genuine content | No         | Protected                     |

### Genuine Content Definition

Content is "genuine" if:

- Created by human (30+ seconds effort)
- Not placeholder/template text
- Not just "work in progress" notice
- Images count for max 15 seconds effort

### Claim Process

1. Submit claim form with reason
2. Wait for review (~3 working days)
3. Decision via email
4. Appeals process available

---

## Summary Checklist

### For MVP Launch

- [ ] Markdown editor with preview
- [ ] Random URL generation
- [ ] Custom URL support
- [ ] Edit code protection
- [ ] Basic markdown rendering
- [ ] Entry deletion
- [ ] Mobile-responsive design

### For Full Feature Parity

- [ ] All MVP features
- [ ] Modify codes
- [ ] Syntax highlighting (200+ languages)
- [ ] Image support with sizing
- [ ] Colored text
- [ ] Underlined text
- [ ] Table of contents
- [ ] Admonitions
- [ ] Spoilers
- [ ] View counter
- [ ] Export (PDF, PNG, JPG, WebP)
- [ ] Passwordless login
- [ ] Metadata system
- [ ] REST API
- [ ] CLI tool
- [ ] URL claim system
- [ ] Rate limiting
- [ ] Content moderation

---

*This specification defines the user-facing features and flows for recreating Rentry.co. Technical implementation details should be determined during development.*
