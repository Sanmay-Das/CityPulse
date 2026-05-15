import { useState, useEffect } from "react";
import { Link } from "react-router-dom";

export default function PrivacyModal() {
  const [visible, setVisible] = useState(false);

  useEffect(() => {
    const accepted = localStorage.getItem("privacy_accepted");
    if (!accepted) setVisible(true);
  }, []);

  const handleAccept = () => {
    localStorage.setItem("privacy_accepted", "true");
    setVisible(false);
  };

  if (!visible) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      {/* Blurred backdrop */}
      <div className="absolute inset-0 bg-slate-950/70 backdrop-blur-sm" />

      {/* Modal */}
      <div className="relative z-10 w-full max-w-2xl mx-4 bg-slate-900 border border-slate-700 rounded-2xl shadow-2xl flex flex-col max-h-[80vh]">
        {/* Header */}
        <div className="px-6 pt-6 pb-4 border-b border-slate-700">
          <div className="flex items-center gap-3 mb-1">
            <div className="w-2 h-6 bg-chicago-red rounded-full" />
            <h2 className="text-white font-bold text-xl">Privacy Policy</h2>
          </div>
          <p className="text-slate-400 text-xs">
            Effective date: May 13, 2026 · citycrimes.io
          </p>
        </div>

        {/* Scrollable content */}
        <div className="overflow-y-auto px-6 py-4 flex-1 text-slate-400 text-sm leading-relaxed space-y-5">
          <Section title="1. Introduction">
            Welcome to citycrimes.io ("we", "us", or "our"). This Privacy Policy explains how we
            collect, use, and protect information when you visit our website. By using
            citycrimes.io, you agree to the terms described in this policy.
          </Section>

          <Section title="2. Information We Collect">
            We may collect the following types of information:
            <ul className="list-disc list-inside mt-2 space-y-1">
              <li>Usage data such as pages visited, time spent on pages, and browser type</li>
              <li>Device information including IP address, operating system, and device identifiers</li>
              <li>Cookie data — small files placed on your device to track browsing behavior</li>
              <li>Natural language queries entered in the Ask AI feature (see Section 3)</li>
            </ul>
            <p className="mt-2">
              We do not directly collect your name, email address, or other personally identifiable
              information unless you voluntarily contact us.
            </p>
          </Section>

          <Section title="3. AI Query Feature (Groq)">
            CityCrimes includes an "Ask AI" feature. Text you enter is transmitted to Groq, Inc.
            via their API for processing. Please do not enter personal or sensitive information in
            the AI query box. Groq's privacy policy is available at{" "}
            <a href="https://groq.com/privacy-policy" target="_blank" rel="noopener noreferrer"
              className="text-chicago-red hover:underline">groq.com/privacy-policy</a>.
          </Section>

          <Section title="4. Google AdSense and Third-Party Advertising">
            We use Google AdSense to display advertisements. Google AdSense uses cookies to serve
            ads based on your prior visits to this and other websites. You may opt out at{" "}
            <a href="https://www.google.com/settings/ads" target="_blank" rel="noopener noreferrer"
              className="text-chicago-red hover:underline">Google Ads Settings</a> or{" "}
            <a href="https://optout.aboutads.info" target="_blank" rel="noopener noreferrer"
              className="text-chicago-red hover:underline">aboutads.info</a>.
            Google's privacy policy is at{" "}
            <a href="https://policies.google.com/privacy" target="_blank" rel="noopener noreferrer"
              className="text-chicago-red hover:underline">policies.google.com/privacy</a>.
          </Section>

          <Section title="5. Cookies">
            Cookies are small text files stored on your device. We and our advertising partners use
            cookies to remember preferences, analyze site usage, and serve relevant ads. You can
            control or disable cookies through your browser settings.
          </Section>

          <Section title="6. How We Use Your Information">
            <ul className="list-disc list-inside mt-2 space-y-1">
              <li>Operate and improve citycrimes.io</li>
              <li>Serve and display advertisements through Google AdSense</li>
              <li>Process natural language queries via the Ask AI feature (Groq)</li>
              <li>Analyze site traffic and usage patterns</li>
              <li>Comply with applicable legal obligations</li>
            </ul>
          </Section>

          <Section title="7. Data Sharing">
            We do not sell your personal information. We may share data with Google LLC (AdSense),
            Groq, Inc. (AI queries), and legal authorities if required by law.
          </Section>

          <Section title="8. GDPR — Rights of EU Users">
            If you are located in the EU, you have the right to access, correct, delete, or
            restrict the processing of your personal data. Contact us at the email below.
          </Section>

          <Section title="9. CCPA — Rights of California Users">
            If you are a California resident, you have the right to know what data we collect,
            request deletion, and opt out of data sales. We do not sell personal data.
          </Section>

          <Section title="10. Children's Privacy">
            citycrimes.io is not directed at children under 13. We do not knowingly collect
            personal information from children.
          </Section>

          <Section title="11. Changes to This Policy">
            We may update this policy from time to time. The effective date at the top reflects the
            most recent revision. Continued use constitutes acceptance of the updated policy.
          </Section>

          <Section title="12. Contact Us">
            Questions? Contact us at{" "}
            <a href="mailto:contact@citycrimes.io" className="text-chicago-red hover:underline">
              contact@citycrimes.io
            </a>{" "}
            or visit{" "}
            <Link to="/privacy" className="text-chicago-red hover:underline">
              citycrimes.io/privacy
            </Link>.
          </Section>
        </div>

        {/* Footer */}
        <div className="px-6 py-4 border-t border-slate-700 flex items-center justify-between gap-4">
          <p className="text-slate-500 text-xs">
            Scroll up to read the full policy before accepting.
          </p>
          <button
            onClick={handleAccept}
            className="bg-chicago-red hover:bg-red-700 text-white text-sm font-semibold px-6 py-2 rounded-lg transition-colors whitespace-nowrap"
          >
            I Understand
          </button>
        </div>
      </div>
    </div>
  );
}

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div>
      <h3 className="text-white font-semibold text-sm mb-1">{title}</h3>
      <div className="text-slate-400 text-sm leading-relaxed">{children}</div>
    </div>
  );
}
