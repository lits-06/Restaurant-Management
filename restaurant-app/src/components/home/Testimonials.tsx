import { useState, useEffect } from 'react';

export default function Testimonials() {
  const [activeIndex, setActiveIndex] = useState(0);

  useEffect(() => {
    const autoScroll = () => {
      if (window.innerWidth >= 768) return; // Không tự scroll trên desktop
      setActiveIndex((prevIndex) => (prevIndex + 1) % 2);
    };

    const interval = setInterval(autoScroll, 5000);
    return () => clearInterval(interval);
  }, []);

  return (
    <section className="py-24 bg-surface">
      <div className="px-margin-desktop w-full max-w-container-max mx-auto text-center mb-16">
        <h2 className="text-3xl font-bold mb-4">The Guest Experience</h2>
        <div className="w-16 h-1 bg-primary-container mx-auto"></div>
      </div>
      <div className="px-margin-desktop w-full max-w-container-max mx-auto overflow-hidden relative">
        <div 
          className="flex gap-8 transition-transform duration-500" 
          style={{ transform: `translateX(-${activeIndex * 100}%)` }}
        >
          {/* Testimonial 1 */}
          <div className="min-w-full md:min-w-[calc(50%-16px)] bg-surface-container p-12 rounded-2xl relative">
            <span className="material-symbols-outlined text-primary-container text-6xl absolute top-6 right-6 opacity-40">format_quote</span>
            <div className="flex items-center gap-4 mb-6">
              <div className="w-12 h-12 rounded-full bg-primary flex items-center justify-center text-white font-bold">JD</div>
              <div>
                <h4 className="font-bold">Julian Devereaux</h4>
                <p className="text-sm opacity-60">Food Critic, Le Guide</p>
              </div>
            </div>
            <p className="italic text-lg leading-relaxed text-on-surface">
              "LuxeBistro doesn't just serve food; they curate moments. The attention to detail in the plating of the wagyu was only matched by the impeccable wine pairing suggestions."
            </p>
          </div>
          {/* Testimonial 2 */}
          <div className="min-w-full md:min-w-[calc(50%-16px)] bg-surface-container-high p-12 rounded-2xl relative">
            <span className="material-symbols-outlined text-primary-container text-6xl absolute top-6 right-6 opacity-40">format_quote</span>
            <div className="flex items-center gap-4 mb-6">
              <div className="w-12 h-12 rounded-full bg-tertiary flex items-center justify-center text-white font-bold">SM</div>
              <div>
                <h4 className="font-bold">Sarah Mitchell</h4>
                <p className="text-sm opacity-60">Private Diner</p>
              </div>
            </div>
            <p className="italic text-lg leading-relaxed text-on-surface">
              "The atmosphere is electric yet refined. We held our anniversary dinner here, and the staff treated us with a level of care I haven't experienced elsewhere in years."
            </p>
          </div>
        </div>
      </div>
    </section>
  );
}