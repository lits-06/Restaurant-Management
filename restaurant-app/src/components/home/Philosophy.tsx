export default function Philosophy() {
  return (
    <section className="py-24 px-margin-desktop w-full max-w-container-max mx-auto">
      <div className="grid grid-cols-12 gap-gutter">
        <div className="col-span-12 md:col-span-5 flex flex-col justify-center space-y-6">
          <h2 className="text-3xl font-bold text-primary">Our Philosophy</h2>
          <p className="text-lg text-on-surface-variant">
            We believe that great food is an act of hospitality that transcends the plate. Our approach is rooted in "Informed Hospitality" — anticipating your needs before they are spoken.
          </p>
          <div className="flex items-center gap-4 py-4 border-l-4 border-primary pl-6 bg-surface-container-low rounded-r-xl">
            <span className="material-symbols-outlined text-primary scale-125">eco</span>
            <div>
              <h4 className="font-bold">Locally Sourced</h4>
              <p className="text-sm opacity-80">Only the finest seasonal ingredients from our network of boutique farms.</p>
            </div>
          </div>
        </div>
        <div className="col-span-12 md:col-span-7 grid grid-cols-2 gap-4">
          <div className="bg-surface-container-high rounded-xl overflow-hidden shadow-sm aspect-square relative group">
            <img 
              className="w-full h-full object-cover transition-transform duration-700 group-hover:scale-110" 
              alt="Professional chef plating micro-greens" 
              src="https://lh3.googleusercontent.com/aida-public/AB6AXuCv6MDbR_0AKkfm66UlBVVU6oYCUhbKKsxhxkhgufS-4qNe0WjXYZgy60N9NPNVlXt266xSqQ4ZPOG-2CRU21cNEZZWV19gjjNRKEhEwlYOd-xK4VOeHR-0iSnFk3GaG3YhzVDuBplM9NVkcuJFDmsTGM-WjOjzpT3o7RLMQxVKjJLEcNmqf1OVbILGib4Fk3iqKwHl4C6IYMVgD2bTdeu475BrMZMtJFT0sgIZZ_qda3e4KiO7LC9v8xnsSqQFkd6nPMSAL-eqls0"
            />
            <div className="absolute inset-0 bg-black/20 opacity-0 group-hover:opacity-100 transition-opacity flex items-end p-6">
              <span className="text-white text-xs font-semibold tracking-widest">PRECISION</span>
            </div>
          </div>
          <div className="bg-surface-container-high rounded-xl overflow-hidden shadow-sm aspect-[4/5] mt-8 relative group">
            <img 
              className="w-full h-full object-cover transition-transform duration-700 group-hover:scale-110" 
              alt="Elegant wine cellar" 
              src="https://lh3.googleusercontent.com/aida-public/AB6AXuBzLLAGhd_cZC3mVJB05iMb5w4cAExZXaPJc3lyfWQ_jEHz-yt5YmWUKb2mF71mwKEeu7YxfQH-P0uk7EN8C7LPGGbPiT-4HVlx9cq6fD5NJpqm5AkeMVqnV8vR2koq0mtC0E2gR7i1smpUiov9z040L5_6imLZ1RoasC5iiszjnMXeu3akEsFeiNXuhn2_l85v2lIRE6boH7cqT94W3iYy55UhB7Wh6kMzL_fdxLCvnXkCicevVTND35qygKatANXBBx6q8Q8RKlU"
            />
            <div className="absolute inset-0 bg-black/20 opacity-0 group-hover:opacity-100 transition-opacity flex items-end p-6">
              <span className="text-white text-xs font-semibold tracking-widest">CURATION</span>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}