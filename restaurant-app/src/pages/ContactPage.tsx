export default function Contact() {
  return (
    <main className="max-w-container-max mx-auto px-margin-mobile md:px-margin-desktop py-16 flex-grow">
      {/* Hero Section */}
      <div className="mb-16 text-center">
        <h1 className="font-headline-xl text-headline-xl text-primary mb-4">Connect With Us</h1>
        <p className="font-body-lg text-body-lg text-on-surface-variant max-w-2xl mx-auto">
          Whether you're planning an intimate dinner or a grand celebration, our concierge team is here to ensure every detail of your LuxeBistro experience is perfection.
        </p>
      </div>

      {/* Main Content Bento Grid */}
      <div className="grid grid-cols-1 gap-gutter max-w-3xl mx-auto">
        {/* Contact Form Card (Trống theo HTML gốc) */}

        {/* Side Info Column */}
        <div className="flex flex-col gap-gutter">
          
          {/* Quick Contact Info */}
          <div className="bg-surface-container-high rounded-xl p-8 shadow-sm">
            <h3 className="font-headline-md text-headline-md mb-6">Concierge Details</h3>
            
            <div className="space-y-6">
              {/* Phone */}
              <div className="flex items-start gap-4">
                <div className="bg-primary-container p-3 rounded-lg text-on-primary-container flex items-center justify-center">
                  <span className="material-symbols-outlined" data-icon="call">call</span>
                </div>
                <div>
                  <p className="font-label-sm text-on-surface-variant">Phone</p>
                  <p className="font-body-lg font-semibold">+1 (555) 890-2345</p>
                </div>
              </div>

              {/* Email */}
              <div className="flex items-start gap-4">
                <div className="bg-primary-container p-3 rounded-lg text-on-primary-container flex items-center justify-center">
                  <span className="material-symbols-outlined" data-icon="mail">mail</span>
                </div>
                <div>
                  <p className="font-label-sm text-on-surface-variant">Email</p>
                  <p className="font-body-lg font-semibold">reservations@luxebistro.com</p>
                </div>
              </div>

              {/* Location */}
              <div className="flex items-start gap-4">
                <div className="bg-primary-container p-3 rounded-lg text-on-primary-container flex items-center justify-center">
                  <span className="material-symbols-outlined" data-icon="location_on">location_on</span>
                </div>
                <div>
                  <p className="font-label-sm text-on-surface-variant">Location</p>
                  <p className="font-body-lg font-semibold">
                    1221 Gastronomy Blvd,<br />New York, NY 10013
                  </p>
                </div>
              </div>
            </div>

            {/* Social Media Links */}
            <div className="mt-8 pt-8 border-t border-outline-variant">
              <p className="font-label-sm text-on-surface-variant mb-4 uppercase tracking-widest text-center">
                Follow the Journey
              </p>
              <div className="flex justify-center gap-4">
                <a className="w-12 h-12 flex items-center justify-center rounded-full bg-surface-container-highest hover:bg-primary hover:text-white transition-all" href="#">
                  <span className="material-symbols-outlined" data-icon="photo_camera">photo_camera</span>
                </a>
                <a className="w-12 h-12 flex items-center justify-center rounded-full bg-surface-container-highest hover:bg-primary hover:text-white transition-all" href="#">
                  <span className="material-symbols-outlined" data-icon="public">public</span>
                </a>
                <a className="w-12 h-12 flex items-center justify-center rounded-full bg-surface-container-highest hover:bg-primary hover:text-white transition-all" href="#">
                  <span className="material-symbols-outlined" data-icon="restaurant">restaurant</span>
                </a>
              </div>
            </div>
          </div>

          {/* Operating Hours */}
          <div className="bg-inverse-surface text-inverse-on-surface rounded-xl p-8 shadow-md">
            <div className="flex items-center gap-3 mb-6">
              <span className="material-symbols-outlined text-primary-fixed" data-icon="schedule">schedule</span>
              <h3 className="font-headline-md text-headline-md">Service Hours</h3>
            </div>
            
            <ul className="space-y-3 font-body-md">
              <li className="flex justify-between border-b border-on-surface-variant pb-2">
                <span>Mon - Thu</span>
                <span className="font-semibold text-primary-fixed">17:00 - 22:30</span>
              </li>
              <li className="flex justify-between border-b border-on-surface-variant pb-2">
                <span>Fri - Sat</span>
                <span className="font-semibold text-primary-fixed">17:00 - 00:00</span>
              </li>
              <li className="flex justify-between border-b border-on-surface-variant pb-2">
                <span>Sunday</span>
                <span className="font-semibold text-primary-fixed">11:00 - 21:00</span>
              </li>
            </ul>
          </div>

        </div>
      </div>
      
      {/* Map Section (Trống theo HTML gốc) */}
    </main>
  );
}