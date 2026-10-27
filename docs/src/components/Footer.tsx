export function Footer() {
  return (
    <footer className="border-t border-zinc-100 py-12 text-xs text-zinc-500">
      <div className="max-w-4xl mx-auto px-6 flex flex-col sm:flex-row items-center justify-between gap-4">
        <div>
          <span className="font-medium text-zinc-800">jvm</span> — Open source under MIT License.
        </div>

        <div className="flex items-center gap-5">
          <a
            href="https://github.com/unstoppable-tanmay/jvm"
            target="_blank"
            rel="noopener noreferrer"
            className="hover:text-zinc-900 transition-colors"
          >
            GitHub
          </a>
          <a
            href="https://www.linkedin.com/in/tanmay-kumar-panda/"
            target="_blank"
            rel="noopener noreferrer"
            className="hover:text-zinc-900 transition-colors"
          >
            LinkedIn
          </a>
          <a
            href="https://github.com/sponsors/unstoppable-tanmay"
            target="_blank"
            rel="noopener noreferrer"
            className="hover:text-zinc-900 transition-colors"
          >
            Sponsor
          </a>
          <a
            href="https://tanmaykumarpanda.in"
            target="_blank"
            rel="noopener noreferrer"
            className="hover:text-zinc-900 transition-colors"
          >
            Portfolio
          </a>
        </div>
      </div>
    </footer>
  );
}
