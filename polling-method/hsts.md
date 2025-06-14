# 📘 HSTS (HTTP Strict Transport Security) — Structured Notes

## 🔐 What is HSTS?

**HSTS (HTTP Strict Transport Security)** is a web security policy mechanism that protects HTTPS websites against downgrade attacks (like SSL stripping) and cookie hijacking.

* Defined in [RFC 6797](https://datatracker.ietf.org/doc/html/rfc6797).

## 📋 How HSTS Works (Browser Behavior)

1. When a user visits a **domain for the first time via HTTPS**, the server can respond with an `Strict-Transport-Security` header.
2. This header typically looks like:

   ```
   Strict-Transport-Security: max-age=31536000; includeSubDomains; preload
   ```
3. The browser stores this HSTS policy locally for the specified `max-age` (in seconds).
4. For all future visits (even if typed as `http://`), the browser will **automatically upgrade the request to HTTPS** — before it hits the network.
5. This is **enforced client-side**.

## 📦 HSTS Preload List

Some domains can be **preloaded into browsers**, meaning the browser knows to enforce HTTPS on them even **before** the first visit.

* This list is **shipped with browsers like Chrome, Firefox, Safari, Edge**, etc.
* To **submit your domain** to the preload list or check if it’s already on it, use this official site:
  🔗 [https://hstspreload.org/?domain=chatgpt.com#submission-form](https://hstspreload.org/?domain=chatgpt.com#submission-form)

## 🌐 chrome://net-internals/#hsts

Chrome allows inspection and manipulation of the local HSTS settings via:
🔗 `chrome://net-internals/#hsts`

Features:

* Add domain to HSTS list manually.
* Query a domain’s HSTS settings.
* Delete domain from HSTS list.

## 🖥 Where Is the HSTS Data Stored?

### 🔹 Chrome:

#### macOS:

* File path: `~/Library/Application Support/Google/Chrome/Default/TransportSecurity`
* You can inspect but **should not manually edit** this binary file.

#### Ubuntu:

* File path: `~/.config/google-chrome/Default/TransportSecurity`

Use `chrome://net-internals/#hsts` for safe manipulation.

### 🔸 Safari (macOS only):

Safari handles HSTS internally but is **not user-configurable** like Chrome.

* **HSTS file location**: Not publicly documented or modifiable.
* **Cannot manually add/delete domains** from Safari’s HSTS list.
* Once HSTS is set by a website in Safari, you must wait until the `max-age` expires or clear all website data.

## ⚠️ Limitations of HSTS

* If a domain enforces HSTS with a long `max-age`, users **cannot revert to HTTP** until expiry.
* Removing your domain from the preload list requires a **browser release cycle**, which can take months.
* Misconfiguring HSTS with preload and long max-age can lead to **loss of access** if HTTPS breaks.
* Not all browsers provide tools to inspect HSTS behavior (e.g., Safari).

---

## ✅ Summary

| Feature                        | Chrome                                        | Safari           |
| ------------------------------ | --------------------------------------------- | ---------------- |
| View HSTS Domains              | `chrome://net-internals/#hsts`                | ❌ Not Available  |
| Add/Delete Domain to HSTS List | ✅ Supported                                   | ❌ Not Supported  |
| File Path (macOS)              | `~/Library/Application Support/...`           | ❌ Undocumented   |
| File Path (Ubuntu)             | `~/.config/google-chrome/...`                 | ❌ N/A            |
| Preload List Submission        | ✅ [hstspreload.org](https://hstspreload.org/) | ❌ Not Applicable |
