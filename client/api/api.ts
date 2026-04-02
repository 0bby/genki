interface getItemsArgs {
  limit: number;
  period: string;
  page: number;
  exercise_id?: number;
}
interface getActivityArgs {
  metric: string;
  step: string;
  range: number;
  month: number;
  year: number;
}
interface timeframe {
  week?: number;
  month?: number;
  year?: number;
  from?: number;
  to?: number;
  period?: string;
}

async function handleJson<T>(r: Response): Promise<T> {
  if (!r.ok) {
    const err = await r.json();
    throw Error(err.error);
  }
  return (await r.json()) as T;
}

async function getWorkouts(
  args: getItemsArgs
): Promise<PaginatedResponse<Workout>> {
  const r = await fetch(
    `/apis/web/v1/workouts?period=${args.period}&limit=${args.limit}&page=${args.page}`
  );
  return handleJson<PaginatedResponse<Workout>>(r);
}

async function getWorkout(id: number): Promise<WorkoutDetail> {
  const r = await fetch(`/apis/web/v1/workout?id=${id}`);
  return handleJson<WorkoutDetail>(r);
}

async function deleteWorkout(id: number): Promise<Response> {
  return fetch(`/apis/web/v1/workout?id=${id}`, { method: "DELETE" });
}

async function getTopExercises(
  args: getItemsArgs
): Promise<PaginatedResponse<Ranked<Exercise>>> {
  const r = await fetch(
    `/apis/web/v1/top-exercises?period=${args.period}&limit=${args.limit}&page=${args.page}`
  );
  return handleJson<PaginatedResponse<Ranked<Exercise>>>(r);
}

async function getTopMuscles(
  args: getItemsArgs
): Promise<PaginatedResponse<Ranked<Muscle>>> {
  const r = await fetch(
    `/apis/web/v1/top-muscles?period=${args.period}&limit=${args.limit}&page=${args.page}`
  );
  return handleJson<PaginatedResponse<Ranked<Muscle>>>(r);
}

async function getActivity(
  args: getActivityArgs
): Promise<ActivityItem[]> {
  const r = await fetch(
    `/apis/web/v1/activity?metric=${args.metric}&step=${args.step}&range=${args.range}&month=${args.month}&year=${args.year}`
  );
  return handleJson<ActivityItem[]>(r);
}

async function getStats(period: string): Promise<FitnessStats> {
  const r = await fetch(`/apis/web/v1/stats?period=${period}`);
  return handleJson<FitnessStats>(r);
}

async function getExercise(id: number): Promise<Exercise> {
  const r = await fetch(`/apis/web/v1/exercise?id=${id}`);
  return handleJson<Exercise>(r);
}

function search(q: string): Promise<SearchResponse> {
  q = encodeURIComponent(q);
  return fetch(`/apis/web/v1/search?q=${q}`).then(
    (r) => r.json() as Promise<SearchResponse>
  );
}

// Health data
async function getSleep(period: string): Promise<SleepLog[]> {
  const r = await fetch(`/apis/web/v1/sleep?period=${period}`);
  return handleJson<SleepLog[]>(r);
}

async function getHeartRate(period: string): Promise<HeartRateDaily[]> {
  const r = await fetch(`/apis/web/v1/heart-rate?period=${period}`);
  return handleJson<HeartRateDaily[]>(r);
}

async function getMeasurements(type: string, period: string): Promise<BodyMeasurement[]> {
  const r = await fetch(`/apis/web/v1/measurements?type=${type}&period=${period}`);
  return handleJson<BodyMeasurement[]>(r);
}

async function getSteps(period: string): Promise<DailySteps[]> {
  const r = await fetch(`/apis/web/v1/steps?period=${period}`);
  return handleJson<DailySteps[]>(r);
}

// Sync
async function getSyncStatus(): Promise<SyncStatus> {
  const r = await fetch(`/apis/web/v1/sync/status`);
  return handleJson<SyncStatus>(r);
}

async function triggerSync(source: string): Promise<Response> {
  return fetch(`/apis/web/v1/sync/trigger?source=${source}`, { method: "POST" });
}

async function initFitbitOAuth(): Promise<{ url: string }> {
  const r = await fetch(`/apis/web/v1/oauth/fitbit/init`, { method: "POST" });
  return handleJson<{ url: string }>(r);
}

// Auth
function login(
  username: string,
  password: string,
  remember: boolean
): Promise<Response> {
  const form = new URLSearchParams();
  form.append("username", username);
  form.append("password", password);
  form.append("remember_me", String(remember));
  return fetch(`/apis/web/v1/login`, {
    method: "POST",
    body: form,
  });
}
function logout(): Promise<Response> {
  return fetch(`/apis/web/v1/logout`, { method: "POST" });
}

function getCfg(): Promise<Config> {
  return fetch(`/apis/web/v1/config`).then((r) => r.json() as Promise<Config>);
}

// User / API keys
function getApiKeys(): Promise<ApiKey[]> {
  return fetch(`/apis/web/v1/user/apikeys`).then(
    (r) => r.json() as Promise<ApiKey[]>
  );
}
const createApiKey = async (label: string): Promise<ApiKey> => {
  const form = new URLSearchParams();
  form.append("label", label);
  const r = await fetch(`/apis/web/v1/user/apikeys`, {
    method: "POST",
    body: form,
  });
  if (!r.ok) {
    let errorMessage = `error: ${r.status}`;
    try {
      const errorData: ApiError = await r.json();
      if (errorData && typeof errorData.error === "string") {
        errorMessage = errorData.error;
      }
    } catch (e) {
      console.error("unexpected api error:", e);
    }
    throw new Error(errorMessage);
  }
  const data: ApiKey = await r.json();
  return data;
};
function deleteApiKey(id: number): Promise<Response> {
  return fetch(`/apis/web/v1/user/apikeys?id=${id}`, { method: "DELETE" });
}
function updateApiKeyLabel(id: number, label: string): Promise<Response> {
  const form = new URLSearchParams();
  form.append("id", String(id));
  form.append("label", label);
  return fetch(`/apis/web/v1/user/apikeys`, {
    method: "PATCH",
    body: form,
  });
}
function updateUser(username: string, password: string) {
  const form = new URLSearchParams();
  form.append("username", username);
  form.append("password", password);
  return fetch(`/apis/web/v1/user`, {
    method: "PATCH",
    body: form,
  });
}

async function getRecapStats(args: timeframe): Promise<RecapStats> {
  const r = await fetch(
    `/apis/web/v1/summary?week=${args.week}&month=${args.month}&year=${args.year}&from=${args.from}&to=${args.to}`
  );
  return handleJson<RecapStats>(r);
}

export {
  getWorkouts,
  getWorkout,
  deleteWorkout,
  getTopExercises,
  getTopMuscles,
  getActivity,
  getStats,
  getExercise,
  search,
  getSleep,
  getHeartRate,
  getMeasurements,
  getSteps,
  getSyncStatus,
  triggerSync,
  initFitbitOAuth,
  login,
  logout,
  getCfg,
  updateUser,
  getApiKeys,
  createApiKey,
  deleteApiKey,
  updateApiKeyLabel,
  getRecapStats,
};

// Types

type Exercise = {
  id: number;
  name: string;
  description?: string;
  category_id?: number;
  category?: ExerciseCategory;
  wger_id?: number;
  muscles?: Muscle[];
  created_at: string;
  total_sets?: number;
  total_reps?: number;
  total_volume_kg?: number;
};

type ExerciseCategory = {
  id: number;
  name: string;
  wger_id?: number;
};

type Muscle = {
  id: number;
  name: string;
  name_en?: string;
  is_front: boolean;
  wger_id?: number;
};

type Workout = {
  id: number;
  user_id: number;
  started_at: string;
  ended_at?: string;
  duration_minutes?: number;
  title?: string;
  notes?: string;
  source: string;
  source_id?: string;
  created_at: string;
};

type WorkoutSet = {
  id: number;
  workout_id: number;
  exercise_id: number;
  exercise?: Exercise;
  set_number: number;
  reps?: number;
  weight_kg?: string;
  duration_seconds?: number;
  rpe?: string;
  logged_at: string;
};

type WorkoutDetail = {
  workout: Workout;
  sets: WorkoutSet[];
};

type SleepLog = {
  id: number;
  user_id: number;
  date: string;
  total_minutes: number;
  deep_minutes?: number;
  light_minutes?: number;
  rem_minutes?: number;
  awake_minutes?: number;
  efficiency?: number;
  start_time?: string;
  end_time?: string;
  source: string;
};

type HeartRateDaily = {
  id: number;
  user_id: number;
  date: string;
  resting_hr?: number;
  avg_hr?: number;
  max_hr?: number;
  source: string;
};

type DailySteps = {
  id: number;
  user_id: number;
  date: string;
  step_count: number;
  source: string;
};

type BodyMeasurement = {
  id: number;
  user_id: number;
  date: string;
  weight_kg?: string;
  body_fat_pct?: string;
  measurement_category?: string;
  measurement_value?: string;
  source: string;
};

type SyncCursor = {
  id: number;
  source: string;
  resource: string;
  last_synced_at: string;
};

type SyncStatus = {
  sources: SyncCursor[];
};

type PaginatedResponse<T> = {
  items: T[];
  total_record_count: number;
  has_next_page: boolean;
  current_page: number;
  items_per_page: number;
};

type Ranked<T> = {
  item: T;
  rank: number;
};

type ActivityItem = {
  start: string;
  value: number;
};

type FitnessStats = {
  workout_count: number;
  exercise_count: number;
  total_active_minutes: number;
  total_steps: number;
  avg_sleep_minutes: number;
  total_sets: number;
  total_reps: number;
  avg_workout_duration: number;
};

type SearchResponse = {
  exercises: Exercise[];
};

type User = {
  id: number;
  username: string;
  role: "user" | "admin";
};

type ApiKey = {
  id: number;
  key: string;
  label: string;
  created_at: Date;
};

type ApiError = {
  error: string;
};

type Config = {
  default_theme: string;
};

type RecapStats = {
  title: string;
  top_exercises: Ranked<Exercise>[];
  top_muscles: Ranked<Muscle>[];
  total_workouts: number;
  total_sets: number;
  total_reps: number;
  total_active_minutes: number;
  avg_workout_duration: number;
  exercises_tried: number;
  new_exercises: number;
  workout_streak: number;
};

export type {
  getItemsArgs,
  getActivityArgs,
  Exercise,
  ExerciseCategory,
  Muscle,
  Workout,
  WorkoutSet,
  WorkoutDetail,
  SleepLog,
  HeartRateDaily,
  DailySteps,
  BodyMeasurement,
  SearchResponse,
  PaginatedResponse,
  Ranked,
  ActivityItem,
  User,
  ApiKey,
  ApiError,
  Config,
  FitnessStats,
  RecapStats,
  SyncStatus,
  SyncCursor,
};
