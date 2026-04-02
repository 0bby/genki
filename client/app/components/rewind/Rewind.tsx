import type { RecapStats } from "api/api";
import RewindStatText from "./RewindStatText";

interface Props {
  stats: RecapStats;
}

export default function Rewind(props: Props) {
  const { stats } = props;

  if (!stats.top_exercises || stats.top_exercises.length === 0) {
    return <p>Not enough data exists to create a Recap for this period.</p>;
  }

  return (
    <div className="flex flex-col gap-7">
      <h2>{stats.title}</h2>

      <div className="flex flex-col gap-2">
        <h4>Top Exercises</h4>
        {stats.top_exercises.map((e, i) => (
          <div key={e.item.id} className="text-sm flex items-center gap-2">
            <span className="font-bold w-6 text-end">{e.rank}</span>
            <span>{e.item.name}</span>
            <span className="text-(--color-fg-tertiary)">
              {e.item.total_sets ?? 0} sets &middot; {e.item.total_reps ?? 0} reps
            </span>
          </div>
        ))}
      </div>

      {stats.top_muscles && stats.top_muscles.length > 0 && (
        <div className="flex flex-col gap-2">
          <h4>Top Muscles</h4>
          {stats.top_muscles.map((m) => (
            <div key={m.item.id} className="text-sm flex items-center gap-2">
              <span className="font-bold w-6 text-end">{m.rank}</span>
              <span>{m.item.name_en || m.item.name}</span>
            </div>
          ))}
        </div>
      )}

      <div className="grid grid-cols-1 sm:grid-cols-3 gap-y-5">
        <RewindStatText figure={`${stats.total_workouts}`} text="Workouts" />
        <RewindStatText figure={`${stats.total_sets}`} text="Total sets" />
        <RewindStatText figure={`${stats.total_reps}`} text="Total reps" />
        <RewindStatText
          figure={`${stats.total_active_minutes}`}
          text="Active minutes"
        />
        <RewindStatText
          figure={`${stats.avg_workout_duration.toFixed(0)}`}
          text="Avg workout (min)"
        />
        <RewindStatText
          figure={`${stats.exercises_tried}`}
          text="Exercises tried"
        />
        <RewindStatText
          figure={`${stats.workout_streak}`}
          text="Best streak (days)"
        />
        <RewindStatText
          figure={`${stats.new_exercises}`}
          text="New exercises"
        />
      </div>
    </div>
  );
}
