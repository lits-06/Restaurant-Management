import { useEffect, useState } from "react";

export function useAutoSlider(length: number) {
  const [index, setIndex] = useState(0);

  useEffect(() => {
    const timer = setInterval(() => {
      setIndex((prev) => (prev + 1) % length);
    }, 5000);

    return () => clearInterval(timer);
  }, [length]);

  return index;
}