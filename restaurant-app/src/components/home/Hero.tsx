export default function Hero() {
  return (
    <section className="relative h-[921px] min-h-[600px] flex items-center overflow-hidden">
      <div className="absolute inset-0 z-0">
        <img 
          className="w-full h-full object-cover" 
          alt="A high-end restaurant dining room" 
          src="https://lh3.googleusercontent.com/aida-public/AB6AXuDsFW58fVAd4Ad_UfpHgXdTDPjqx67LYIyuiZkU9JgR-HnEsF3ZjE_xSV5oh5NQ8oXaBMhKfT1P7X_1oiQs3ZeqMAt6TSRMICTbdhYeAW3NRmnBKJaIAxYntQbLhUeuO3Ef3_jxikLNOrAMPfsHDJ71JH_O7Z4aYL5freojvOzHmJgP3-ai23mluQgqx6Qge648S-F4Q9IVDppma2doCvsoD7TjWJRc5TDMsmBt-ZUsEgWQ-0Z8sSgtqPjCMcAu1pRv5pD3OCr4E54"
        />
        <div className="absolute inset-0 hero-gradient"></div>
      </div>
      <div className="relative z-10 px-margin-desktop w-full max-w-container-max mx-auto">
        <div className="max-w-2xl space-y-6">
          <span className="text-primary-fixed text-xs font-semibold tracking-[0.2em] uppercase">
            Est. 1998 — Paris &amp; New York
          </span>
          <h1 className="text-5xl md:text-6xl text-surface-container-lowest leading-tight font-bold">
            Culinary Artistry <br />
            <span className="text-primary-fixed-dim italic font-normal">Defined by Precision.</span>
          </h1>
          <p className="text-surface-variant text-lg max-w-lg">
            Experience the harmonious blend of traditional technique and modern innovation in every dish we serve.
          </p>
          <div className="flex gap-4 pt-4">
            <button className="bg-primary-container text-on-primary-container px-8 py-4 rounded-lg text-xl font-bold flex items-center gap-2 hover:shadow-lg transition-all active:scale-95">
              Reserve Your Table
              <span className="material-symbols-outlined">arrow_forward</span>
            </button>
          </div>
        </div>
      </div>
    </section>
  );
}