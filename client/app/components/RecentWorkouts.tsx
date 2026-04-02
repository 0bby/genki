import { useQuery } from "@tanstack/react-query";
import { timeSince } from "~/utils/utils";
import { getWorkouts, type getItemsArgs, type Workout } from "api/api";
import { Link } from "react-router";

interface Props {
  limit: number;
}

export default function RecentWorkouts(props: Props) {
  const { isPending, isError, data, error } = useQuery({
    queryKey: [
      "recent-workouts",
      { limit: props.limit, period: "all_time", page: 0 },
    ],
    queryFn: ({ queryKey }) => getWorkouts(queryKey[1] as getItemsArgs),
  });

  const header = "Recent workouts";

  if (isPending) {
    return (
      <div className="w-[300px] sm:w-[500px]">
        <h3>{header}</h3>
        <p>Loading...</p>
      </div>
    );
  } else if (isError) {
    return (
      <div className="w-[300px] sm:w-[500px]">
        <h3>{header}</h3>
        <p className="error">Error: {error.message}</p>
      </div>
    );
  }

  return (
    <div className="text-sm sm:text-[16px]">
      <h3 className="hover:underline">
        <Link to="/workouts">{header}</Link>
      </h3>
      <table className="-ml-4">
        <tbody>
          {data.items.map((workout) => (
            <tr
              key={`workout_${workout.id}`}
              className="group hover:bg-[--color-bg-secondary]"
            >
              <td
                className="color-fg-tertiary pr-2 sm:pr-4 text-sm whitespace-nowrap w-0"
                title={new Date(workout.started_at).toString()}
              >
                {timeSince(new Date(workout.started_at))}
              </td>
              <td className="text-ellipsis overflow-hidden max-w-[400px] sm:max-w-[600px]">
                <Link
                  className="hover:text-[--color-fg-secondary]"
                  to={`/workout/${workout.id}`}
                >
                  {workout.title || "Workout"}
                </Link>
              </td>
              <td className="color-fg-tertiary text-sm pl-2">
                {workout.duration_minutes
                  ? `${workout.duration_minutes} min`
                  : ""}
              </td>
              <td className="color-fg-tertiary text-xs pl-2">{workout.source}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
