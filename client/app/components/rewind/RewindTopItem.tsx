import type { Ranked } from "api/api";

type TopItemProps<T> = {
  title: string;
  items: Ranked<T>[];
  getLabel: (item: T) => string;
};

export function RewindTopItem<
  T extends { id: string | number }
>({ title, items, getLabel }: TopItemProps<T>) {
  const [top, ...rest] = items;

  if (!top) return null;

  return (
    <div className="flex flex-col gap-1">
      <h4 className="-mb-1">{title}</h4>
      <h2>{getLabel(top.item)}</h2>
      {rest.map((e) => (
        <div key={e.item.id} className="text-sm">
          <span className="font-bold mr-1">{e.rank}.</span>
          {getLabel(e.item)}
        </div>
      ))}
    </div>
  );
}
