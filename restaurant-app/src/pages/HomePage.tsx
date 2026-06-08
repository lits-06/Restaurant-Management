import Header from "@/components/layout/Header";
import Footer from "@/components/layout/Footer";

import HeroSection from "@/components/home/Hero";
import PhilosophySection from "@/components/home/Philosophy";
import FeaturedDishesSection from "@/components/home/FeaturedDishes";
import TestimonialSection from "@/components/home/Testimonials";
import CTASection from "@/components/home/CTA";

export default function HomePage() {
  return (
    <>
        <HeroSection />
        <PhilosophySection />
        <FeaturedDishesSection />
        <TestimonialSection />
        <CTASection />
    </>
  );
}