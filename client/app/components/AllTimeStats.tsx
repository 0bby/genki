import { useQuery } from "@tanstack/react-query";
import { getStats, type FitnessStats } from "api/api";

export default function AllTimeStats() {
  const { isPending, isError, data, error } = useQuery({
    queryKey: ["stats", "all_time"],
    queryFn: ({ queryKey }) => getStats(queryKey[1]),
  });

  const header = "All time stats";

  if (isPending) {
    return (
      <div>
        <h3>{header}</h3>
        <p>Loading...</p>
      </div>
    );
  } else if (isError) {
    return (
      <div>
        <h3>{header}</h3>
        <p className="error">Error: {error.message}</p>
      </div>
    );
  }

  const n = "header-font font-bold text-xl";

  return (
    <div>
      <h3>{header}</h3>
      <div>
        <span className={n}>{data.workout_count}</span> Workouts
      </div>
      <div>
        <span className={n}>{data.exercise_count}</span> Exercises
      </div>
      <div>
        <span className={n}>{data.total_sets}</span> Sets
      </div>
      <div>
        <span className={n}>{data.total_reps}</span> Reps
      </div>
      <div>
        <span
          className={n}
          title={Math.floor(data.total_active_minutes / 60) + " hours"}
        >
          {data.total_active_minutes}
        </span>{" "}
        Active Minutes
      </div>
      <div>
        <span className={n}>{data.total_steps.toLocaleString()}</span> Steps
      </div>
    </div>
  );
}
