import { Link } from "react-router-dom";

export default function PrivacyPolicy() {
  return (
    <div className="h-full overflow-y-auto bg-slate-950 text-slate-300">
      {/* Header */}
      <header className="bg-slate-900/90 border-b border-slate-700 px-6 py-3 flex items-center gap-3">
        <div className="w-2 h-6 bg-chicago-red rounded-full" />
        <Link to="/" className="text-white font-bold text-lg tracking-tight hover:text-slate-200">
          CityCrimes
        </Link>
        <span className="text-slate-500 text-sm">/ Privacy Policy</span>
      </header>

      {/* Content */}
      <main className="max-w-3xl mx-auto px-6 py-12">
        <h1 className="text-white text-3xl font-bold mb-2">Privacy Policy</h1>
        <p className="text-slate-500 text-sm mb-10">
          Effective date: May 13, 2026 · Website: citycrimes.io
        </p>

        <Section title="1. Introduction">
          <p>
            Welcome to citycrimes.io ("we", "us", or "our"). This Privacy Policy explains how we collect,
            use, and protect information when you visit our website. By using citycrimes.io, you agree to
            the terms described in this policy.
          </p>
        </Section>

        <Section title="2. Information We Collect">
          <p className="mb-3">We may collect the following types of information:</p>
          <ul className="list-disc list-inside space-y-1 text-slate-400">
            <li>Usage data such as pages visited, time spent on pages, and browser type</li>
            <li>Device information including IP address, operating system, and device identifiers</li>
            <li>Cookie data — small files placed on your device to track browsing behavior</li>
            <li>Natural language queries entered in the Ask AI feature (see Section 3)</li>
          </ul>
          <p className="mt-3">
            We do not directly collect your name, email address, or other personally identifiable
            information unless you voluntarily contact us.
          </p>
        </Section>

        <Section title="3. AI Query Feature (Groq)">
          <p>
            CityCrimes includes an "Ask AI" feature that allows you to ask natural language questions
            about crime data. Text you enter in this feature is transmitted to Groq, Inc. via their
            API for processing. Groq's privacy policy governs how they handle this data and is
            available at{" "}
            <a
              href="https://groq.com/privacy-policy"
              target="_blank"
              rel="noopener noreferrer"
              className="text-chicago-red hover:underline"
            >
              groq.com/privacy-policy
            </a>
            . Please do not enter personal or sensitive information in the AI query box.
          </p>
        </Section>

        <Section title="4. Google AdSense and Third-Party Advertising">
          <p className="mb-3">
            We use Google AdSense to display advertisements on citycrimes.io. Google AdSense uses
            cookies and similar tracking technologies to serve ads based on your prior visits to this
            and other websites.
          </p>
          <ul className="list-disc list-inside space-y-1 text-slate-400 mb-3">
            <li>
              Third-party vendors, including Google, use cookies to serve ads based on a user's prior
              visits to citycrimes.io or other websites on the internet.
            </li>
            <li>
              Google's use of advertising cookies enables it and its partners to serve ads to you
              based on your visit to our site and/or other sites on the internet.
            </li>
            <li>
              You may opt out of personalized advertising by visiting{" "}
              <a
                href="https://www.google.com/settings/ads"
                target="_blank"
                rel="noopener noreferrer"
                className="text-chicago-red hover:underline"
              >
                Google Ads Settings
              </a>
              .
            </li>
            <li>
              You may also opt out of third-party vendor cookies for interest-based advertising by
              visiting{" "}
              <a
                href="https://optout.aboutads.info"
                target="_blank"
                rel="noopener noreferrer"
                className="text-chicago-red hover:underline"
              >
                aboutads.info
              </a>
              .
            </li>
          </ul>
          <p>
            Google's privacy policy governing its use of data is available at{" "}
            <a
              href="https://policies.google.com/privacy"
              target="_blank"
              rel="noopener noreferrer"
              className="text-chicago-red hover:underline"
            >
              policies.google.com/privacy
            </a>
            .
          </p>
        </Section>

        <Section title="5. Cookies">
          <p className="mb-3">
            Cookies are small text files stored on your device when you visit a website. We and our
            advertising partners use cookies to:
          </p>
          <ul className="list-disc list-inside space-y-1 text-slate-400 mb-3">
            <li>Remember your preferences and settings</li>
            <li>Analyze how our site is used</li>
            <li>Serve relevant advertisements</li>
          </ul>
          <p>
            You can control or disable cookies through your browser settings. Note that disabling
            cookies may affect your experience on citycrimes.io and other websites.
          </p>
        </Section>

        <Section title="6. How We Use Your Information">
          <p className="mb-3">We use the information collected to:</p>
          <ul className="list-disc list-inside space-y-1 text-slate-400">
            <li>Operate and improve citycrimes.io</li>
            <li>Serve and display advertisements through Google AdSense</li>
            <li>Process natural language queries via the Ask AI feature (Groq)</li>
            <li>Analyze site traffic and usage patterns</li>
            <li>Comply with applicable legal obligations</li>
          </ul>
        </Section>

        <Section title="7. Data Sharing">
          <p className="mb-3">
            We do not sell your personal information. We may share data with:
          </p>
          <ul className="list-disc list-inside space-y-1 text-slate-400">
            <li>Google LLC, for the purpose of serving ads via AdSense</li>
            <li>Groq, Inc., for processing natural language queries entered in the Ask AI feature</li>
            <li>Legal authorities, if required by law</li>
          </ul>
        </Section>

        <Section title="8. GDPR — Rights of EU Users">
          <p>
            If you are located in the European Union, you have the right to access, correct, delete,
            or restrict the processing of your personal data. You also have the right to object to
            processing and to data portability. To exercise these rights, contact us at the email below.
          </p>
        </Section>

        <Section title="9. CCPA — Rights of California Users">
          <p>
            If you are a California resident, you have the right to know what personal data we
            collect, to request deletion of your data, and to opt out of the sale of personal data.
            We do not sell personal data. To exercise your rights, contact us at the email below.
          </p>
        </Section>

        <Section title="10. Children's Privacy">
          <p>
            citycrimes.io is not directed at children under the age of 13. We do not knowingly
            collect personal information from children. If you believe a child has provided us with
            personal information, please contact us and we will delete it.
          </p>
        </Section>

        <Section title="11. Changes to This Policy">
          <p>
            We may update this Privacy Policy from time to time. The effective date at the top of
            this page will reflect the most recent revision. Continued use of citycrimes.io after
            changes constitutes acceptance of the updated policy.
          </p>
        </Section>

        <Section title="12. Contact Us">
          <p>
            If you have questions about this Privacy Policy, please contact us at:
          </p>
          <p className="mt-2 text-slate-400">
            Website:{" "}
            <a href="https://citycrimes.io" className="text-chicago-red hover:underline">
              citycrimes.io
            </a>
            <br />
            Email:{" "}
            <a href="mailto:contact@citycrimes.io" className="text-chicago-red hover:underline">
              contact@citycrimes.io
            </a>
          </p>
        </Section>

        <div className="mt-12 pt-6 border-t border-slate-800">
          <Link
            to="/"
            className="text-slate-400 hover:text-white text-sm transition-colors"
          >
            ← Back to CityCrimes
          </Link>
        </div>
      </main>
    </div>
  );
}

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <section className="mb-8">
      <h2 className="text-white font-semibold text-lg mb-3 pb-1 border-b border-slate-800">
        {title}
      </h2>
      <div className="text-slate-400 text-sm leading-relaxed space-y-2">{children}</div>
    </section>
  );
}
