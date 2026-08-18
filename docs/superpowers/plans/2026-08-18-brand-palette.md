# Brand Color Palette (Navy/Teal/Orange) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the de-facto indigo/gray brand colors in `frontend/tenant` with a new Navy/Teal/Orange palette, without touching semantic status colors (success/warning/danger/info) or neutral background/border/secondary-text colors.

**Architecture:** Add three custom color scales (`navy`, `teal` override, `orange` override) via Tailwind v4's `@theme` block, add a matching PrimeVue `definePreset` (primary → teal) so built-in components pick it up automatically, then migrate navigation/shared components and each `views/modules/*` directory in separate reviewable batches.

**Tech Stack:** Vue 3, Tailwind CSS v4 (`@theme`, no `tailwind.config.js`), PrimeVue 4 (Aura preset + `definePreset`), no frontend test runner in this repo (verification is `npm run build` + manual visual smoke-check, matching project convention).

**Spec:** `docs/superpowers/specs/2026-08-18-brand-palette-design.md`

## Global Constraints

- Scope is `frontend/tenant` only. `frontend/platform-admin` is out of scope.
- Status/semantic colors (`emerald`=success, `amber`=warning, `rose`=danger, `sky`/`blue`=info) are never touched.
- Neutral colors (`gray`/`slate` for background, card, border, secondary text) are never touched.
- Only **light-mode** `text-gray-800`/`text-gray-900` (primary/heading text) is remapped to navy. The paired `dark:text-gray-100` (or similar light-on-dark) stays as-is — do not touch dark-mode text color pairings.
- `indigo-*` utility classes (link, active icon/background, focus ring, info badge) become `teal-*` with the same shade number and same `dark:`/opacity modifiers (e.g. `text-indigo-600 dark:text-indigo-400` → `text-teal-600 dark:text-teal-400`).
- Orange is applied narrowly and explicitly per element (primary CTA buttons, notification badges/highlights) — never as a blanket replace of another color family.
- No frontend test suite exists in this repo — "tests" for this plan means `npm run build` succeeding plus a manual visual check in the dev server (light + dark mode).

---

## File Structure

- Modify: `frontend/tenant/src/assets/styles/main.css` — add `@theme` block with `navy`/`teal`/`orange` scales.
- Modify: `frontend/tenant/src/main.js` — add `definePreset` extending Aura, primary → teal.
- Modify: `frontend/tenant/src/components/AppLayout.vue` (and any sidebar/nav sub-component it renders) — navy for nav background/text/active-item.
- Modify: `frontend/tenant/src/components/**/*.vue` (shared components used across modules) — indigo → teal, light-mode `text-gray-800` → `text-navy-800`.
- Modify: `frontend/tenant/src/views/modules/{attendance,leave,payroll,performance,recruitment,reimbursement,training,employee,employeemovement,organization,job,jobvalues,approval,notification,competency,career-intelligence,workforce-intelligence}/**/*.vue` — same rule, one directory per task/commit.

---

## Task 1: Add Tailwind `@theme` color tokens

**Files:**
- Modify: `frontend/tenant/src/assets/styles/main.css:1-10`

**Interfaces:**
- Produces: Tailwind utility classes `bg-navy-{50..950}`, `text-navy-{50..950}`, `border-navy-{50..950}` (new), and overridden `teal-{50..950}` / `orange-{50..950}` scales usable by every later task exactly like any other Tailwind color family (e.g. `bg-teal-50 dark:bg-teal-500/10 text-teal-600 dark:text-teal-400`).

- [ ] **Step 1: Add the `@theme` block**

Insert immediately after the `@import "tailwindcss";` line (line 1) in `frontend/tenant/src/assets/styles/main.css`:

```css
@theme {
  /* Navy — primary/heading text, navigation. Base #1B2A41 sits at navy-800
     (matches how it replaces text-gray-800 for primary text). */
  --color-navy-50: #f6f6f7;
  --color-navy-100: #eaecee;
  --color-navy-200: #d1d4d9;
  --color-navy-300: #b6bbc2;
  --color-navy-400: #969da8;
  --color-navy-500: #767f8d;
  --color-navy-600: #566172;
  --color-navy-700: #364458;
  --color-navy-800: #1b2a41;
  --color-navy-900: #162234;
  --color-navy-950: #121b2a;

  /* Teal override — brand secondary: link, active icon/badge, info badge.
     Base #1B7F93 sits at teal-600 (matches how it replaces text-indigo-600). */
  --color-teal-50: #f4f9fa;
  --color-teal-100: #e6f1f3;
  --color-teal-200: #c8e0e5;
  --color-teal-300: #a8ced6;
  --color-teal-400: #84bac5;
  --color-teal-500: #549fae;
  --color-teal-600: #1b7f93;
  --color-teal-700: #176c7d;
  --color-teal-800: #135967;
  --color-teal-900: #0f4651;
  --color-teal-950: #0b333b;

  /* Orange override — CTA / notification / highlight, used sparingly.
     Base #F5941E sits at orange-500 (matches typical CTA button shade). */
  --color-orange-50: #fffaf4;
  --color-orange-100: #fef3e6;
  --color-orange-200: #fde4c7;
  --color-orange-300: #fbd4a5;
  --color-orange-400: #f9b96d;
  --color-orange-500: #f5941e;
  --color-orange-600: #d8821a;
  --color-orange-700: #b86f17;
  --color-orange-800: #935912;
  --color-orange-900: #6e430e;
  --color-orange-950: #4a2c09;
}
```

- [ ] **Step 2: Verify Tailwind picks up the new tokens**

Run: `cd frontend/tenant && npm run build`
Expected: build succeeds with no CSS errors (Tailwind v4 fails the build on malformed `@theme` syntax, so a clean build is sufficient verification here — there is no unit test for CSS tokens).

- [ ] **Step 3: Commit**

```bash
git add frontend/tenant/src/assets/styles/main.css
git commit -m "feat(theme): tambah token warna navy + override teal/orange untuk brand palette baru"
```

---

## Task 2: Add PrimeVue `definePreset` (primary → teal)

**Files:**
- Modify: `frontend/tenant/src/main.js` (wherever `Aura` is currently imported and passed to `PrimeVue`)

**Interfaces:**
- Consumes: nothing from Task 1 directly (PrimeVue preset uses its own hex tokens, not Tailwind's `@theme` — they must be kept visually consistent by using the same hex values).
- Produces: PrimeVue's `primary` semantic color token now resolves to the teal scale, affecting default `Button`, `InputText` focus ring, `Menu`/`PanelMenu` active item, `Checkbox`, `RadioButton`, etc. app-wide with no further per-component change needed.

- [ ] **Step 1: Read current PrimeVue setup**

Open `frontend/tenant/src/main.js` and locate the existing:
```js
import Aura from '@primevue/themes/aura'
...
app.use(PrimeVue, {
  theme: {
    preset: Aura,
    options: { darkModeSelector: '.p-dark' }
  }
})
```

- [ ] **Step 2: Define and wire the teal preset**

Add near the top of `frontend/tenant/src/main.js` (after the `Aura` import):

```js
import { definePreset } from '@primevue/themes'

const TealPreset = definePreset(Aura, {
  semantic: {
    primary: {
      50: '{teal.50}',
      100: '{teal.100}',
      200: '{teal.200}',
      300: '{teal.300}',
      400: '{teal.400}',
      500: '{teal.500}',
      600: '#1b7f93',
      700: '#176c7d',
      800: '#135967',
      900: '#0f4651',
      950: '#0b333b'
    }
  }
})
```

Then change the `PrimeVue` registration to use `TealPreset` instead of `Aura`:

```js
app.use(PrimeVue, {
  theme: {
    preset: TealPreset,
    options: { darkModeSelector: '.p-dark' }
  }
})
```

(If `main.js` already defines `50`–`500` for a `teal` alias elsewhere, reuse it — otherwise inline the same hex values used in Task 1's `--color-teal-*` tokens for `50`–`500` so both systems match: `50:#f4f9fa 100:#e6f1f3 200:#c8e0e5 300:#a8ced6 400:#84bac5 500:#549fae`.)

- [ ] **Step 3: Build and smoke-check**

Run: `cd frontend/tenant && npm run build`
Expected: build succeeds.

Run: `cd frontend/tenant && npm run dev`, open the app in a browser, confirm:
- Any default PrimeVue `Button` (not `severity="secondary"`/etc.) now renders teal, not the old default blue.
- Focus ring on a text input is teal.
- Toggle dark mode — teal primary still visible and readable (not washed out).

- [ ] **Step 4: Commit**

```bash
git add frontend/tenant/src/main.js
git commit -m "feat(theme): definePreset PrimeVue primary=teal untuk komponen bawaan"
```

---

## Task 3: Migrate navigation (AppLayout / Sidebar) to Navy

**Files:**
- Modify: `frontend/tenant/src/components/AppLayout.vue`
- Modify: any sidebar/nav sub-component `AppLayout.vue` renders (identify via its `<script>` imports in Step 1)

**Interfaces:**
- Consumes: `navy-*` Tailwind classes from Task 1.
- Produces: nothing consumed by later tasks — this is a leaf UI change.

- [ ] **Step 1: Identify all navigation-related files**

Run:
```bash
grep -n "^import" frontend/tenant/src/components/AppLayout.vue | grep -i "sidebar\|nav\|menu"
```
Note every component path returned — each is in scope for this task alongside `AppLayout.vue` itself.

- [ ] **Step 2: Find current color classes in scope**

Run:
```bash
grep -n "indigo-\|gray-800\|gray-900\|bg-white dark:bg-gray" frontend/tenant/src/components/AppLayout.vue
```
(Repeat for each file found in Step 1.)

- [ ] **Step 3: Apply the mapping**

For every match from Step 2, in `AppLayout.vue` and its nav sub-components:
- Sidebar background (light mode): if currently a plain `bg-white`/`bg-gray-50`, change the sidebar's own light-mode background to `bg-navy-800` and make its text/icons light (`text-navy-50`/`text-white`) so the nav reads as a dark navy bar — dark mode is unaffected (nav can keep its existing dark-mode `bg-gray-900`/`bg-gray-800` since it's already dark).
- Active/selected nav item: `bg-indigo-50 text-indigo-600` (or similar) → `bg-teal-600/10 text-teal-300` (teal accent against the navy background) — do NOT use `text-navy-*` for the active item, navy is the background here, teal marks "active" per the spec's mapping.
- Any hover state using `indigo-` → `teal-`.
- Any `text-gray-800`/`text-gray-900` used for nav labels (light mode) → `text-navy-50` (since the nav background is now dark navy, labels need to be light, not `text-navy-800` — this is the one place the "text-gray-800 → text-navy-800" default rule does NOT apply, because the surface underneath changed from light to dark).

- [ ] **Step 4: Build and visually verify**

Run: `cd frontend/tenant && npm run build`
Expected: build succeeds.

Run dev server, confirm: sidebar renders navy background, active nav item has a visible teal accent, hover states are teal, text is legible in both light and dark page mode (the nav's own dark navy look should not change when toggling app dark mode — only the rest of the page does).

- [ ] **Step 5: Commit**

```bash
git add frontend/tenant/src/components/AppLayout.vue
git commit -m "feat(theme): migrasi navigasi (AppLayout/Sidebar) ke Navy + accent Teal"
```

---

## Task 4: Migrate shared components (`frontend/tenant/src/components/**`)

**Files:**
- Modify: all `.vue` files under `frontend/tenant/src/components/` **except** the ones already handled in Task 3 (nav/sidebar).

**Interfaces:**
- Consumes: `navy-*`/`teal-*`/`orange-*` classes from Task 1.
- Produces: nothing consumed by later tasks.

- [ ] **Step 1: List remaining files with brand-color usage**

Run:
```bash
grep -rl "indigo-\|text-gray-800\b" frontend/tenant/src/components --include="*.vue" | grep -v -i "sidebar\|applayout"
```
This is the exact file list for this task.

- [ ] **Step 2: For each file in the list, apply the mapping**

For every occurrence:
- `text-indigo-{N} dark:text-indigo-{M}` → `text-teal-{N} dark:text-teal-{M}` (keep the same numbers, only swap the family name).
- `bg-indigo-{N} dark:bg-indigo-{M}/{opacity}` → `bg-teal-{N} dark:bg-teal-{M}/{opacity}`.
- `border-indigo-{N}` (hover/focus borders, e.g. `hover:border-indigo-300`) → `border-teal-{N}`.
- `ring-indigo-{N}` / `focus-visible:ring-indigo-{N}/{opacity}` → `ring-teal-{N}` / `focus-visible:ring-teal-{N}/{opacity}`.
- `text-gray-800` (light-mode only — leave any `dark:text-gray-*` pair untouched) → `text-navy-800`.
- Any component whose sole purpose is a primary call-to-action button (e.g. a shared `<PrimaryButton>`/`<ConfirmDeleteDialog>` confirm action, if such a shared component exists) — check its label/intent; if it represents the single main affirmative action of a dialog, its background may use `bg-orange-500 hover:bg-orange-600` instead of teal. If unsure whether a shared button component is a "CTA" vs a generic action button, leave it as teal (per the spec: orange is applied narrowly and explicitly, default to the more conservative choice).

- [ ] **Step 3: Build and verify**

Run: `cd frontend/tenant && npm run build`
Expected: build succeeds, no unused-file or syntax errors.

- [ ] **Step 4: Commit**

```bash
git add frontend/tenant/src/components
git commit -m "feat(theme): migrasi shared components ke palet Navy/Teal/Orange"
```

---

## Task 5: Migrate `views/modules/attendance/**`

**Files:**
- Modify: all `.vue` files under `frontend/tenant/src/views/modules/attendance/`

**Interfaces:**
- Consumes: `navy-*`/`teal-*`/`orange-*` classes from Task 1.

- [ ] **Step 1: List files with brand-color usage in this module**

Run:
```bash
grep -rl "indigo-\|text-gray-800\b" frontend/tenant/src/views/modules/attendance --include="*.vue"
```

- [ ] **Step 2: Apply the same mapping as Task 4 Step 2** to every file in the list (indigo → teal per shade/modifier, light-mode `text-gray-800` → `text-navy-800`, orange only for genuine primary CTA buttons and left as teal by default when unsure).

- [ ] **Step 3: Build and verify**

Run: `cd frontend/tenant && npm run build`
Expected: build succeeds.

Spot-check in dev server: open the Attendance hub page and at least one sub-page (e.g. `AttendanceOvertime.vue` or `BusinessTravelList.vue`), confirm links/active icons are teal and headings are navy in light mode, and dark mode still looks correct (unchanged from before, since only light-mode gray-800 was touched).

- [ ] **Step 4: Commit**

```bash
git add frontend/tenant/src/views/modules/attendance
git commit -m "feat(theme): migrasi module Attendance ke palet Navy/Teal/Orange"
```

---

## Task 6: Migrate `views/modules/leave/**`

**Files:** all `.vue` under `frontend/tenant/src/views/modules/leave/`

- [ ] **Step 1:** `grep -rl "indigo-\|text-gray-800\b" frontend/tenant/src/views/modules/leave --include="*.vue"`
- [ ] **Step 2:** For every occurrence in the files listed, apply:
  - `text-indigo-{N} dark:text-indigo-{M}` -> `text-teal-{N} dark:text-teal-{M}` (keep the same numbers, only swap the family name).
  - `bg-indigo-{N} dark:bg-indigo-{M}/{opacity}` -> `bg-teal-{N} dark:bg-teal-{M}/{opacity}`.
  - `border-indigo-{N}` (e.g. `hover:border-indigo-300`) -> `border-teal-{N}`.
  - `ring-indigo-{N}` / `focus-visible:ring-indigo-{N}/{opacity}` -> `ring-teal-{N}` / `focus-visible:ring-teal-{N}/{opacity}`.
  - `text-gray-800` (light-mode only -- leave any paired `dark:text-gray-*` untouched) -> `text-navy-800`.
  - A shared button/element that is the single main affirmative CTA of its dialog/page may use `bg-orange-500 hover:bg-orange-600` instead of teal; if unsure, default to teal.
- [ ] **Step 3:** `cd frontend/tenant && npm run build` — expect success; spot-check the Leave list/detail page in dev server.
- [ ] **Step 4:** `git add frontend/tenant/src/views/modules/leave && git commit -m "feat(theme): migrasi module Leave ke palet Navy/Teal/Orange"`

---

## Task 7: Migrate `views/modules/payroll/**`

**Files:** all `.vue` under `frontend/tenant/src/views/modules/payroll/`

- [ ] **Step 1:** `grep -rl "indigo-\|text-gray-800\b" frontend/tenant/src/views/modules/payroll --include="*.vue"`
- [ ] **Step 2:** For every occurrence in the files listed, apply:
  - `text-indigo-{N} dark:text-indigo-{M}` -> `text-teal-{N} dark:text-teal-{M}` (keep the same numbers, only swap the family name).
  - `bg-indigo-{N} dark:bg-indigo-{M}/{opacity}` -> `bg-teal-{N} dark:bg-teal-{M}/{opacity}`.
  - `border-indigo-{N}` (e.g. `hover:border-indigo-300`) -> `border-teal-{N}`.
  - `ring-indigo-{N}` / `focus-visible:ring-indigo-{N}/{opacity}` -> `ring-teal-{N}` / `focus-visible:ring-teal-{N}/{opacity}`.
  - `text-gray-800` (light-mode only -- leave any paired `dark:text-gray-*` untouched) -> `text-navy-800`.
  - A shared button/element that is the single main affirmative CTA of its dialog/page may use `bg-orange-500 hover:bg-orange-600` instead of teal; if unsure, default to teal.
- [ ] **Step 3:** `cd frontend/tenant && npm run build` — expect success; spot-check a Payroll run detail page in dev server.
- [ ] **Step 4:** `git add frontend/tenant/src/views/modules/payroll && git commit -m "feat(theme): migrasi module Payroll ke palet Navy/Teal/Orange"`

---

## Task 8: Migrate `views/modules/performance/**`

**Files:** all `.vue` under `frontend/tenant/src/views/modules/performance/`

- [ ] **Step 1:** `grep -rl "indigo-\|text-gray-800\b" frontend/tenant/src/views/modules/performance --include="*.vue"`
- [ ] **Step 2:** For every occurrence in the files listed, apply:
  - `text-indigo-{N} dark:text-indigo-{M}` -> `text-teal-{N} dark:text-teal-{M}` (keep the same numbers, only swap the family name).
  - `bg-indigo-{N} dark:bg-indigo-{M}/{opacity}` -> `bg-teal-{N} dark:bg-teal-{M}/{opacity}`.
  - `border-indigo-{N}` (e.g. `hover:border-indigo-300`) -> `border-teal-{N}`.
  - `ring-indigo-{N}` / `focus-visible:ring-indigo-{N}/{opacity}` -> `ring-teal-{N}` / `focus-visible:ring-teal-{N}/{opacity}`.
  - `text-gray-800` (light-mode only -- leave any paired `dark:text-gray-*` untouched) -> `text-navy-800`.
  - A shared button/element that is the single main affirmative CTA of its dialog/page may use `bg-orange-500 hover:bg-orange-600` instead of teal; if unsure, default to teal.
- [ ] **Step 3:** `cd frontend/tenant && npm run build` — expect success; spot-check KPI/OKR index page in dev server.
- [ ] **Step 4:** `git add frontend/tenant/src/views/modules/performance && git commit -m "feat(theme): migrasi module Performance ke palet Navy/Teal/Orange"`

---

## Task 9: Migrate `views/modules/recruitment/**`

**Files:** all `.vue` under `frontend/tenant/src/views/modules/recruitment/`

- [ ] **Step 1:** `grep -rl "indigo-\|text-gray-800\b" frontend/tenant/src/views/modules/recruitment --include="*.vue"`
- [ ] **Step 2:** For every occurrence in the files listed, apply:
  - `text-indigo-{N} dark:text-indigo-{M}` -> `text-teal-{N} dark:text-teal-{M}` (keep the same numbers, only swap the family name).
  - `bg-indigo-{N} dark:bg-indigo-{M}/{opacity}` -> `bg-teal-{N} dark:bg-teal-{M}/{opacity}`.
  - `border-indigo-{N}` (e.g. `hover:border-indigo-300`) -> `border-teal-{N}`.
  - `ring-indigo-{N}` / `focus-visible:ring-indigo-{N}/{opacity}` -> `ring-teal-{N}` / `focus-visible:ring-teal-{N}/{opacity}`.
  - `text-gray-800` (light-mode only -- leave any paired `dark:text-gray-*` untouched) -> `text-navy-800`.
  - A shared button/element that is the single main affirmative CTA of its dialog/page may use `bg-orange-500 hover:bg-orange-600` instead of teal; if unsure, default to teal.
- [ ] **Step 3:** `cd frontend/tenant && npm run build` — expect success; spot-check Requisitions/Candidates page in dev server.
- [ ] **Step 4:** `git add frontend/tenant/src/views/modules/recruitment && git commit -m "feat(theme): migrasi module Recruitment ke palet Navy/Teal/Orange"`

---

## Task 10: Migrate `views/modules/reimbursement/**`

**Files:** all `.vue` under `frontend/tenant/src/views/modules/reimbursement/`

- [ ] **Step 1:** `grep -rl "indigo-\|text-gray-800\b" frontend/tenant/src/views/modules/reimbursement --include="*.vue"`
- [ ] **Step 2:** For every occurrence in the files listed, apply:
  - `text-indigo-{N} dark:text-indigo-{M}` -> `text-teal-{N} dark:text-teal-{M}` (keep the same numbers, only swap the family name).
  - `bg-indigo-{N} dark:bg-indigo-{M}/{opacity}` -> `bg-teal-{N} dark:bg-teal-{M}/{opacity}`.
  - `border-indigo-{N}` (e.g. `hover:border-indigo-300`) -> `border-teal-{N}`.
  - `ring-indigo-{N}` / `focus-visible:ring-indigo-{N}/{opacity}` -> `ring-teal-{N}` / `focus-visible:ring-teal-{N}/{opacity}`.
  - `text-gray-800` (light-mode only -- leave any paired `dark:text-gray-*` untouched) -> `text-navy-800`.
  - A shared button/element that is the single main affirmative CTA of its dialog/page may use `bg-orange-500 hover:bg-orange-600` instead of teal; if unsure, default to teal.
- [ ] **Step 3:** `cd frontend/tenant && npm run build` — expect success; spot-check Reimbursement Requests page in dev server.
- [ ] **Step 4:** `git add frontend/tenant/src/views/modules/reimbursement && git commit -m "feat(theme): migrasi module Reimbursement ke palet Navy/Teal/Orange"`

---

## Task 11: Migrate `views/modules/training/**`

**Files:** all `.vue` under `frontend/tenant/src/views/modules/training/`

- [ ] **Step 1:** `grep -rl "indigo-\|text-gray-800\b" frontend/tenant/src/views/modules/training --include="*.vue"`
- [ ] **Step 2:** For every occurrence in the files listed, apply:
  - `text-indigo-{N} dark:text-indigo-{M}` -> `text-teal-{N} dark:text-teal-{M}` (keep the same numbers, only swap the family name).
  - `bg-indigo-{N} dark:bg-indigo-{M}/{opacity}` -> `bg-teal-{N} dark:bg-teal-{M}/{opacity}`.
  - `border-indigo-{N}` (e.g. `hover:border-indigo-300`) -> `border-teal-{N}`.
  - `ring-indigo-{N}` / `focus-visible:ring-indigo-{N}/{opacity}` -> `ring-teal-{N}` / `focus-visible:ring-teal-{N}/{opacity}`.
  - `text-gray-800` (light-mode only -- leave any paired `dark:text-gray-*` untouched) -> `text-navy-800`.
  - A shared button/element that is the single main affirmative CTA of its dialog/page may use `bg-orange-500 hover:bg-orange-600` instead of teal; if unsure, default to teal.
- [ ] **Step 3:** `cd frontend/tenant && npm run build` — expect success; spot-check Training Courses/Sessions page in dev server.
- [ ] **Step 4:** `git add frontend/tenant/src/views/modules/training && git commit -m "feat(theme): migrasi module Training ke palet Navy/Teal/Orange"`

---

## Task 12: Migrate `views/modules/employee/**`

**Files:** all `.vue` under `frontend/tenant/src/views/modules/employee/`

- [ ] **Step 1:** `grep -rl "indigo-\|text-gray-800\b" frontend/tenant/src/views/modules/employee --include="*.vue"`
- [ ] **Step 2:** For every occurrence in the files listed, apply:
  - `text-indigo-{N} dark:text-indigo-{M}` -> `text-teal-{N} dark:text-teal-{M}` (keep the same numbers, only swap the family name).
  - `bg-indigo-{N} dark:bg-indigo-{M}/{opacity}` -> `bg-teal-{N} dark:bg-teal-{M}/{opacity}`.
  - `border-indigo-{N}` (e.g. `hover:border-indigo-300`) -> `border-teal-{N}`.
  - `ring-indigo-{N}` / `focus-visible:ring-indigo-{N}/{opacity}` -> `ring-teal-{N}` / `focus-visible:ring-teal-{N}/{opacity}`.
  - `text-gray-800` (light-mode only -- leave any paired `dark:text-gray-*` untouched) -> `text-navy-800`.
  - A shared button/element that is the single main affirmative CTA of its dialog/page may use `bg-orange-500 hover:bg-orange-600` instead of teal; if unsure, default to teal.
- [ ] **Step 3:** `cd frontend/tenant && npm run build` — expect success; spot-check Employee list/detail page in dev server.
- [ ] **Step 4:** `git add frontend/tenant/src/views/modules/employee && git commit -m "feat(theme): migrasi module Employee ke palet Navy/Teal/Orange"`

---

## Task 13: Migrate `views/modules/employeemovement/**`

**Files:** all `.vue` under `frontend/tenant/src/views/modules/employeemovement/`

- [ ] **Step 1:** `grep -rl "indigo-\|text-gray-800\b" frontend/tenant/src/views/modules/employeemovement --include="*.vue"`
- [ ] **Step 2:** For every occurrence in the files listed, apply:
  - `text-indigo-{N} dark:text-indigo-{M}` -> `text-teal-{N} dark:text-teal-{M}` (keep the same numbers, only swap the family name).
  - `bg-indigo-{N} dark:bg-indigo-{M}/{opacity}` -> `bg-teal-{N} dark:bg-teal-{M}/{opacity}`.
  - `border-indigo-{N}` (e.g. `hover:border-indigo-300`) -> `border-teal-{N}`.
  - `ring-indigo-{N}` / `focus-visible:ring-indigo-{N}/{opacity}` -> `ring-teal-{N}` / `focus-visible:ring-teal-{N}/{opacity}`.
  - `text-gray-800` (light-mode only -- leave any paired `dark:text-gray-*` untouched) -> `text-navy-800`.
  - A shared button/element that is the single main affirmative CTA of its dialog/page may use `bg-orange-500 hover:bg-orange-600` instead of teal; if unsure, default to teal.
- [ ] **Step 3:** `cd frontend/tenant && npm run build` — expect success; spot-check Employee Movements page in dev server.
- [ ] **Step 4:** `git add frontend/tenant/src/views/modules/employeemovement && git commit -m "feat(theme): migrasi module Employee Movement ke palet Navy/Teal/Orange"`

---

## Task 14: Migrate `views/modules/organization/**`

**Files:** all `.vue` under `frontend/tenant/src/views/modules/organization/`

- [ ] **Step 1:** `grep -rl "indigo-\|text-gray-800\b" frontend/tenant/src/views/modules/organization --include="*.vue"`
- [ ] **Step 2:** For every occurrence in the files listed, apply:
  - `text-indigo-{N} dark:text-indigo-{M}` -> `text-teal-{N} dark:text-teal-{M}` (keep the same numbers, only swap the family name).
  - `bg-indigo-{N} dark:bg-indigo-{M}/{opacity}` -> `bg-teal-{N} dark:bg-teal-{M}/{opacity}`.
  - `border-indigo-{N}` (e.g. `hover:border-indigo-300`) -> `border-teal-{N}`.
  - `ring-indigo-{N}` / `focus-visible:ring-indigo-{N}/{opacity}` -> `ring-teal-{N}` / `focus-visible:ring-teal-{N}/{opacity}`.
  - `text-gray-800` (light-mode only -- leave any paired `dark:text-gray-*` untouched) -> `text-navy-800`.
  - A shared button/element that is the single main affirmative CTA of its dialog/page may use `bg-orange-500 hover:bg-orange-600` instead of teal; if unsure, default to teal.
- [ ] **Step 3:** `cd frontend/tenant && npm run build` — expect success; spot-check Organizations page in dev server.
- [ ] **Step 4:** `git add frontend/tenant/src/views/modules/organization && git commit -m "feat(theme): migrasi module Organization ke palet Navy/Teal/Orange"`

---

## Task 15: Migrate `views/modules/job/**` and `views/modules/jobvalues/**`

**Files:** all `.vue` under `frontend/tenant/src/views/modules/job/` and `frontend/tenant/src/views/modules/jobvalues/`

- [ ] **Step 1:** `grep -rl "indigo-\|text-gray-800\b" frontend/tenant/src/views/modules/job frontend/tenant/src/views/modules/jobvalues --include="*.vue"`
- [ ] **Step 2:** For every occurrence in the files listed, apply:
  - `text-indigo-{N} dark:text-indigo-{M}` -> `text-teal-{N} dark:text-teal-{M}` (keep the same numbers, only swap the family name).
  - `bg-indigo-{N} dark:bg-indigo-{M}/{opacity}` -> `bg-teal-{N} dark:bg-teal-{M}/{opacity}`.
  - `border-indigo-{N}` (e.g. `hover:border-indigo-300`) -> `border-teal-{N}`.
  - `ring-indigo-{N}` / `focus-visible:ring-indigo-{N}/{opacity}` -> `ring-teal-{N}` / `focus-visible:ring-teal-{N}/{opacity}`.
  - `text-gray-800` (light-mode only -- leave any paired `dark:text-gray-*` untouched) -> `text-navy-800`.
  - A shared button/element that is the single main affirmative CTA of its dialog/page may use `bg-orange-500 hover:bg-orange-600` instead of teal; if unsure, default to teal.
- [ ] **Step 3:** `cd frontend/tenant && npm run build` — expect success; spot-check Job Management page in dev server.
- [ ] **Step 4:** `git add frontend/tenant/src/views/modules/job frontend/tenant/src/views/modules/jobvalues && git commit -m "feat(theme): migrasi module Job/Job Values ke palet Navy/Teal/Orange"`

---

## Task 16: Migrate `views/modules/approval/**`

**Files:** all `.vue` under `frontend/tenant/src/views/modules/approval/`

- [ ] **Step 1:** `grep -rl "indigo-\|text-gray-800\b" frontend/tenant/src/views/modules/approval --include="*.vue"`
- [ ] **Step 2:** For every occurrence in the files listed, apply:
  - `text-indigo-{N} dark:text-indigo-{M}` -> `text-teal-{N} dark:text-teal-{M}` (keep the same numbers, only swap the family name).
  - `bg-indigo-{N} dark:bg-indigo-{M}/{opacity}` -> `bg-teal-{N} dark:bg-teal-{M}/{opacity}`.
  - `border-indigo-{N}` (e.g. `hover:border-indigo-300`) -> `border-teal-{N}`.
  - `ring-indigo-{N}` / `focus-visible:ring-indigo-{N}/{opacity}` -> `ring-teal-{N}` / `focus-visible:ring-teal-{N}/{opacity}`.
  - `text-gray-800` (light-mode only -- leave any paired `dark:text-gray-*` untouched) -> `text-navy-800`.
  - A shared button/element that is the single main affirmative CTA of its dialog/page may use `bg-orange-500 hover:bg-orange-600` instead of teal; if unsure, default to teal.
- [ ] **Step 3:** `cd frontend/tenant && npm run build` — expect success; spot-check Approval Flows/Approvals page in dev server.
- [ ] **Step 4:** `git add frontend/tenant/src/views/modules/approval && git commit -m "feat(theme): migrasi module Approval ke palet Navy/Teal/Orange"`

---

## Task 17: Migrate `views/modules/notification/**`

**Files:** all `.vue` under `frontend/tenant/src/views/modules/notification/`

- [ ] **Step 1:** `grep -rl "indigo-\|text-gray-800\b" frontend/tenant/src/views/modules/notification --include="*.vue"`
- [ ] **Step 2:** For every occurrence in the files listed, apply:
  - `text-indigo-{N} dark:text-indigo-{M}` -> `text-teal-{N} dark:text-teal-{M}` (keep the same numbers, only swap the family name).
  - `bg-indigo-{N} dark:bg-indigo-{M}/{opacity}` -> `bg-teal-{N} dark:bg-teal-{M}/{opacity}`.
  - `border-indigo-{N}` (e.g. `hover:border-indigo-300`) -> `border-teal-{N}`.
  - `ring-indigo-{N}` / `focus-visible:ring-indigo-{N}/{opacity}` -> `ring-teal-{N}` / `focus-visible:ring-teal-{N}/{opacity}`.
  - `text-gray-800` (light-mode only -- leave any paired `dark:text-gray-*` untouched) -> `text-navy-800`.
  - A shared button/element that is the single main affirmative CTA of its dialog/page may use `bg-orange-500 hover:bg-orange-600` instead of teal; if unsure, default to teal. This module is a good candidate for the orange "notification highlight" rule — check unread/badge indicators specifically: if an existing unread-count badge or highlighted row currently uses `indigo-`/`blue-` to mean "unread"/"new", change it to `orange-500`/`orange-600` instead of teal (this is the "notification" case the spec calls out for orange, not a generic link/icon).
- [ ] **Step 3:** `cd frontend/tenant && npm run build` — expect success; spot-check Notifications page in dev server, confirm unread indicator now reads as orange.
- [ ] **Step 4:** `git add frontend/tenant/src/views/modules/notification && git commit -m "feat(theme): migrasi module Notification ke palet Navy/Teal/Orange (unread badge -> orange)"`

---

## Task 18: Migrate `views/modules/competency/**`

**Files:** all `.vue` under `frontend/tenant/src/views/modules/competency/`

- [ ] **Step 1:** `grep -rl "indigo-\|text-gray-800\b" frontend/tenant/src/views/modules/competency --include="*.vue"`
- [ ] **Step 2:** For every occurrence in the files listed, apply:
  - `text-indigo-{N} dark:text-indigo-{M}` -> `text-teal-{N} dark:text-teal-{M}` (keep the same numbers, only swap the family name).
  - `bg-indigo-{N} dark:bg-indigo-{M}/{opacity}` -> `bg-teal-{N} dark:bg-teal-{M}/{opacity}`.
  - `border-indigo-{N}` (e.g. `hover:border-indigo-300`) -> `border-teal-{N}`.
  - `ring-indigo-{N}` / `focus-visible:ring-indigo-{N}/{opacity}` -> `ring-teal-{N}` / `focus-visible:ring-teal-{N}/{opacity}`.
  - `text-gray-800` (light-mode only -- leave any paired `dark:text-gray-*` untouched) -> `text-navy-800`.
  - A shared button/element that is the single main affirmative CTA of its dialog/page may use `bg-orange-500 hover:bg-orange-600` instead of teal; if unsure, default to teal.
- [ ] **Step 3:** `cd frontend/tenant && npm run build` — expect success; spot-check Competency 360 Assessment Result page in dev server.
- [ ] **Step 4:** `git add frontend/tenant/src/views/modules/competency && git commit -m "feat(theme): migrasi module Competency 360 ke palet Navy/Teal/Orange"`

---

## Task 19: Migrate `views/modules/career-intelligence/**`

**Files:** all `.vue` under `frontend/tenant/src/views/modules/career-intelligence/`

- [ ] **Step 1:** `grep -rl "indigo-\|text-gray-800\b" frontend/tenant/src/views/modules/career-intelligence --include="*.vue"`
- [ ] **Step 2:** For every occurrence in the files listed, apply:
  - `text-indigo-{N} dark:text-indigo-{M}` -> `text-teal-{N} dark:text-teal-{M}` (keep the same numbers, only swap the family name).
  - `bg-indigo-{N} dark:bg-indigo-{M}/{opacity}` -> `bg-teal-{N} dark:bg-teal-{M}/{opacity}`.
  - `border-indigo-{N}` (e.g. `hover:border-indigo-300`) -> `border-teal-{N}`.
  - `ring-indigo-{N}` / `focus-visible:ring-indigo-{N}/{opacity}` -> `ring-teal-{N}` / `focus-visible:ring-teal-{N}/{opacity}`.
  - `text-gray-800` (light-mode only -- leave any paired `dark:text-gray-*` untouched) -> `text-navy-800`.
  - A shared button/element that is the single main affirmative CTA of its dialog/page may use `bg-orange-500 hover:bg-orange-600` instead of teal; if unsure, default to teal.
- [ ] **Step 3:** `cd frontend/tenant && npm run build` — expect success; spot-check Career Paths/Intelligence page in dev server.
- [ ] **Step 4:** `git add frontend/tenant/src/views/modules/career-intelligence && git commit -m "feat(theme): migrasi module Career Intelligence ke palet Navy/Teal/Orange"`

---

## Task 20: Migrate `views/modules/workforce-intelligence/**`

**Files:** all `.vue` under `frontend/tenant/src/views/modules/workforce-intelligence/`

- [ ] **Step 1:** `grep -rl "indigo-\|text-gray-800\b" frontend/tenant/src/views/modules/workforce-intelligence --include="*.vue"`
- [ ] **Step 2:** For every occurrence in the files listed, apply:
  - `text-indigo-{N} dark:text-indigo-{M}` -> `text-teal-{N} dark:text-teal-{M}` (keep the same numbers, only swap the family name).
  - `bg-indigo-{N} dark:bg-indigo-{M}/{opacity}` -> `bg-teal-{N} dark:bg-teal-{M}/{opacity}`.
  - `border-indigo-{N}` (e.g. `hover:border-indigo-300`) -> `border-teal-{N}`.
  - `ring-indigo-{N}` / `focus-visible:ring-indigo-{N}/{opacity}` -> `ring-teal-{N}` / `focus-visible:ring-teal-{N}/{opacity}`.
  - `text-gray-800` (light-mode only -- leave any paired `dark:text-gray-*` untouched) -> `text-navy-800`.
  - A shared button/element that is the single main affirmative CTA of its dialog/page may use `bg-orange-500 hover:bg-orange-600` instead of teal; if unsure, default to teal.
- [ ] **Step 3:** `cd frontend/tenant && npm run build` — expect success; spot-check Workforce Intelligence dashboard in dev server.
- [ ] **Step 4:** `git add frontend/tenant/src/views/modules/workforce-intelligence && git commit -m "feat(theme): migrasi module Workforce Intelligence ke palet Navy/Teal/Orange"`

---

## Task 21: Final sweep — dashboards, views root, and dark-mode CSS overrides

**Files:**
- Modify: `frontend/tenant/src/components/dashboard/**/*.vue` (KPI/dashboard widgets — distinct directory from Task 4's shared components, e.g. `EmploymentDashboard.vue`, `HRAttendanceLeaveDashboard.vue`).
- Modify: any `.vue` directly under `frontend/tenant/src/views/` (not inside `modules/`, e.g. `Login.vue`, `Dashboard.vue`, `Profile.vue`).
- Modify: `frontend/tenant/src/assets/styles/main.css:83-87` (`.badge-draft` uses `bg-indigo-50 text-indigo-700 border border-indigo-200` — this is a shared badge class, not covered by any per-file grep above).

**Interfaces:**
- Consumes: `navy-*`/`teal-*` classes from Task 1.

- [ ] **Step 1: Sweep remaining indigo usage across the whole app**

Run:
```bash
grep -rl "indigo-" frontend/tenant/src --include="*.vue"
grep -rn "indigo-" frontend/tenant/src/assets/styles/main.css
```
This should now return only files not covered by Tasks 3–20 (dashboards, top-level views, and the CSS badge class). Apply the Task 4 Step 2 mapping to every `.vue` match.

- [ ] **Step 2: Fix the `.badge-draft` CSS class**

In `frontend/tenant/src/assets/styles/main.css`, change:
```css
.badge-draft {
  @apply bg-indigo-50 text-indigo-700 border border-indigo-200 rounded-md px-2 py-0.5 text-xs font-medium;
}
```
to:
```css
.badge-draft {
  @apply bg-teal-50 text-teal-700 border border-teal-200 rounded-md px-2 py-0.5 text-xs font-medium;
}
```
("Draft" is a link/info-adjacent badge state, not a status color per the spec's kept list of success/warning/danger/info — it maps to teal like the rest of the indigo family.)

- [ ] **Step 3: Confirm zero remaining indigo usage**

Run:
```bash
grep -rl "indigo-" frontend/tenant/src --include="*.vue" --include="*.css"
```
Expected: no output (empty result).

- [ ] **Step 4: Final full build + dark mode smoke test**

Run: `cd frontend/tenant && npm run build`
Expected: build succeeds.

Run dev server, walk through: Login page, main Dashboard, one dashboard widget (e.g. `HRAttendanceLeaveDashboard.vue`), and toggle dark mode on/off — confirm navy/teal/orange render correctly and no leftover indigo/blue is visible where teal or navy was expected.

- [ ] **Step 5: Commit**

```bash
git add frontend/tenant/src/components/dashboard frontend/tenant/src/views frontend/tenant/src/assets/styles/main.css
git commit -m "feat(theme): sweep akhir indigo -> teal (dashboard, top-level views, badge-draft) + verifikasi nihil indigo tersisa"
```
