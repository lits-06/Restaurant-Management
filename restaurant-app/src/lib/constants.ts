import {
  hero,
  salmon,
  wagyu,
  ravioli,
} from "@/assets/images";

export const DISHES = [
  {
    id: 1,
    name: "Atlantic Glazed Salmon",
    description:
      "Wild-caught, infused with ginger-soy and served over a bed of truffled leak purée.",
    price: 42,
    image: salmon,
  },
  {
    id: 2,
    name: "Heritage Wagyu A5",
    description:
      "Slow-roasted beef with rosemary jus and garlic-confit fingerling potatoes.",
    price: 125,
    image: wagyu,
  },
  {
    id: 3,
    name: "Truffle Agnolotti",
    description:
      "Handmade pasta filled with ricotta and sage, finished with fresh Umbrian truffles.",
    price: 38,
    image: ravioli,
  },
];

export const TESTIMONIALS = [
  {
    id: 1,
    name: "Julian Devereaux",
    role: "Food Critic",
    content:
      "LuxeBistro doesn't just serve food; they curate moments.",
    avatar: "JD",
  },
  {
    id: 2,
    name: "Sarah Mitchell",
    role: "Private Diner",
    content:
      "The atmosphere is electric yet refined.",
    avatar: "SM",
  },
];