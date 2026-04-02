import { type SearchResponse } from "api/api";
import { useState } from "react";
import SearchResultItem from "./SearchResultItem";
import SearchResultSelectorItem from "./SearchResultSelectorItem";

interface Props {
  data?: SearchResponse;
  onSelect: Function;
  selectorMode?: boolean;
}

export default function SearchResults({ data, onSelect, selectorMode }: Props) {
  const [selected, setSelected] = useState(0);
  const classes = "flex flex-col items-start bg rounded w-full";
  const hClasses = "pt-4 pb-2";

  const selectItem = (title: string, id: number) => {
    if (selected === id) {
      setSelected(0);
      onSelect({ id: 0, title: "" });
    } else {
      setSelected(id);
      onSelect({ id: id, title: title });
    }
  };

  if (!data) {
    return <></>;
  }

  return (
    <div className="w-full">
      {data.exercises && data.exercises.length > 0 && (
        <>
          <h3 className={hClasses}>Exercises</h3>
          <div className={classes}>
            {data.exercises.map((exercise) =>
              selectorMode ? (
                <SearchResultSelectorItem
                  key={exercise.id}
                  id={exercise.id}
                  onClick={() => selectItem(exercise.name, exercise.id)}
                  text={exercise.name}
                  subtext={exercise.category?.name}
                  active={selected === exercise.id}
                />
              ) : (
                <SearchResultItem
                  key={exercise.id}
                  to={`/exercise/${exercise.id}`}
                  onClick={() => onSelect(exercise.id)}
                  text={exercise.name}
                  subtext={exercise.category?.name}
                />
              )
            )}
          </div>
        </>
      )}
    </div>
  );
}
