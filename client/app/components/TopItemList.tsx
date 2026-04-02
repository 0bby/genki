import { Link } from "react-router";
import {
  type Exercise,
  type Muscle,
  type PaginatedResponse,
  type Ranked,
} from "api/api";

type Item = Exercise | Muscle;

interface Props<T extends Ranked<Item>> {
  data: PaginatedResponse<T>;
  separators?: ConstrainBoolean;
  ranked?: boolean;
  type: "exercise" | "muscle";
  className?: string;
}

export default function TopItemList<T extends Ranked<Item>>({
  data,
  separators,
  type,
  className,
  ranked,
}: Props<T>) {
  return (
    <div className={`flex flex-col gap-1 ${className} min-w-[200px]`}>
      {data.items.map((item, index) => {
        const key = `${type}-${item.item.id}`;
        return (
          <div
            key={key}
            style={{ fontSize: 12 }}
            className={`${
              separators && index !== data.items.length - 1
                ? "border-b border-(--color-fg-tertiary) mb-1 pb-2"
                : ""
            }`}
          >
            <ItemCard
              ranked={ranked}
              rank={item.rank}
              item={item.item}
              type={type}
            />
          </div>
        );
      })}
    </div>
  );
}

function ItemCard({
  item,
  type,
  rank,
  ranked,
}: {
  item: Item;
  type: "exercise" | "muscle";
  rank: number;
  ranked?: boolean;
}) {
  const itemClasses = "flex items-center gap-2";

  switch (type) {
    case "exercise": {
      const exercise = item as Exercise;
      return (
        <div style={{ fontSize: 12 }} className={itemClasses}>
          {ranked && <div className="w-7 text-end">{rank}</div>}
          <div className="min-w-[48px] w-[48px] h-[48px] rounded bg-(--color-bg-tertiary) flex items-center justify-center text-lg">
            🏋️
          </div>
          <div>
            <Link
              to={`/exercise/${exercise.id}`}
              className="hover:text-(--color-fg-secondary)"
            >
              <span style={{ fontSize: 14 }}>{exercise.name}</span>
            </Link>
            {exercise.category && (
              <div className="color-fg-secondary">{exercise.category.name}</div>
            )}
            <div className="color-fg-secondary">
              {exercise.total_sets ?? 0} sets &middot; {exercise.total_reps ?? 0} reps
            </div>
          </div>
        </div>
      );
    }
    case "muscle": {
      const muscle = item as Muscle;
      return (
        <div style={{ fontSize: 12 }} className={itemClasses}>
          {ranked && <div className="w-7 text-end">{rank}</div>}
          <div className="min-w-[48px] w-[48px] h-[48px] rounded bg-(--color-bg-tertiary) flex items-center justify-center text-lg">
            💪
          </div>
          <div>
            <span style={{ fontSize: 14 }}>{muscle.name_en || muscle.name}</span>
            <div className="color-fg-secondary">
              {muscle.is_front ? "Front" : "Back"}
            </div>
          </div>
        </div>
      );
    }
  }
}
