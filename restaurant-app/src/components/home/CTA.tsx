export default function CTA() {
  return (
    <section className="py-24 px-margin-desktop">
      <div className="w-full max-w-container-max mx-auto bg-inverse-surface text-inverse-on-surface rounded-[2rem] overflow-hidden flex flex-col md:flex-row">
        <div className="p-12 md:p-24 flex-1 space-y-8">
          <h2 className="text-5xl font-bold">Begin Your Journey.</h2>
          <p className="text-surface-variant text-lg opacity-80">
            Join us for an evening of unparalleled flavor and hospitality. Reservations are recommended at least 48 hours in advance.
          </p>
          <div className="flex flex-wrap gap-4">
            <button className="bg-primary text-on-primary px-10 py-5 rounded-lg text-xl font-bold hover:scale-105 transition-all shadow-lg">
              Secure a Table
            </button>
            <button className="border border-outline text-surface-container-lowest px-10 py-5 rounded-lg text-xl font-bold hover:bg-surface-variant/10 transition-all">
              Private Events
            </button>
          </div>
        </div>
        <div className="hidden md:block w-1/3 relative">
          <img 
            className="w-full h-full object-cover" 
            alt="Sophisticated cocktail preparation" 
            src="https://lh3.googleusercontent.com/aida-public/AB6AXuClOHv59KP8BoOwjGjfRlLSgvakjtkxHsvdlrxJBQ1T60SdSG3Cmw0U4S1JGRi3kJQXGOHkiAP_7Jl-sjMvK7efNOI2wj4c2LXk2hBKgTPJ_9oD6hVE3F89JsUe_I-yjEkj8mz7Aso4xT48AFfeRhH0xtswZZSrIsPKBD2OuHjf1uVD9aEm7bC1-YMkAFIXeAmixuJx2mofIlsDbcKsnN-RAnaBVZMhIZX9MZUlNC-4gQzLQygC-an0M4LnGW1JHEoPMyadz34P5jA"
          />
        </div>
      </div>
    </section>
  );
}