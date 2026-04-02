import { type RouteConfig, index, route } from "@react-router/dev/routes";

export default [
  index("routes/Home.tsx"),
  route("/exercise/:id", "routes/ExerciseDetail.tsx"),
  route("/workout/:id", "routes/WorkoutDetail.tsx"),
  route("/chart/top-exercises", "routes/Charts/ExerciseChart.tsx"),
  route("/chart/top-muscles", "routes/Charts/MuscleChart.tsx"),
  route("/workouts", "routes/Charts/WorkoutList.tsx"),
  route("/recap", "routes/RecapPage.tsx"),
  route("/health", "routes/HealthPage.tsx"),
  route("/theme-helper", "routes/ThemeHelper.tsx"),
] satisfies RouteConfig;
