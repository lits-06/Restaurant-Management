import { Dish } from "@/lib/types";

interface Props {
  dish: Dish;
}

export default function DishCard({
  dish,
}: Props) {
  return (
    <div className="bg-white rounded-xl shadow">
      <img
        src={dish.image}
        alt={dish.name}
        className="aspect-[4/3] w-full object-cover"
      />

      <div className="p-4">
        <h3 className="font-semibold text-xl">
          {dish.name}
        </h3>

        <p className="text-gray-500">
          {dish.description}
        </p>

        <p className="mt-4 font-bold">
          ${dish.price}
        </p>
      </div>
    </div>
  );
}