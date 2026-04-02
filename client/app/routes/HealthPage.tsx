import { useQuery } from "@tanstack/react-query";
import { getSleep, getSteps, getHeartRate } from "api/api";
import { useState } from "react";
import PeriodSelector from "~/components/PeriodSelector";

export default function HealthPage() {
  const [period, setPeriod] = useState("week");

  const { data: sleepData } = useQuery({
    queryKey: ["sleep", period],
    queryFn: () => getSleep(period),
  });

  const { data: stepsData } = useQuery({
    queryKey: ["steps", period],
    queryFn: () => getSteps(period),
  });

  const { data: hrData } = useQuery({
    queryKey: ["heart-rate", period],
    queryFn: () => getHeartRate(period),
  });

  const pgTitle = "Health - Genki";

  return (
    <div className="w-full min-h-screen">
      <title>{pgTitle}</title>
      <div className="w-19/20 sm:w-17/20 mx-auto pt-6 sm:pt-12">
        <h1>Health</h1>
        <PeriodSelector current={period} setter={setPeriod} disableCache />

        <div className="mt-10 flex flex-col gap-10">
          {/* Steps */}
          <div>
            <h3>Steps</h3>
            {stepsData && stepsData.length > 0 ? (
              <table className="text-sm mt-2">
                <thead>
                  <tr className="text-(--color-fg-tertiary)">
                    <th className="pr-6 text-left">Date</th>
                    <th className="pr-6 text-left">Steps</th>
                  </tr>
                </thead>
                <tbody>
                  {stepsData.map((s) => (
                    <tr key={s.id}>
                      <td className="pr-6">
                        {new Date(s.date).toLocaleDateString()}
                      </td>
                      <td className="pr-6 font-bold">
                        {s.step_count.toLocaleString()}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            ) : (
              <p className="text-sm text-(--color-fg-tertiary)">No step data</p>
            )}
          </div>

          {/* Sleep */}
          <div>
            <h3>Sleep</h3>
            {sleepData && sleepData.length > 0 ? (
              <table className="text-sm mt-2">
                <thead>
                  <tr className="text-(--color-fg-tertiary)">
                    <th className="pr-6 text-left">Date</th>
                    <th className="pr-6 text-left">Total</th>
                    <th className="pr-6 text-left">Deep</th>
                    <th className="pr-6 text-left">REM</th>
                    <th className="pr-6 text-left">Efficiency</th>
                  </tr>
                </thead>
                <tbody>
                  {sleepData.map((s) => (
                    <tr key={s.id}>
                      <td className="pr-6">
                        {new Date(s.date).toLocaleDateString()}
                      </td>
                      <td className="pr-6">
                        {(s.total_minutes / 60).toFixed(1)}h
                      </td>
                      <td className="pr-6">
                        {s.deep_minutes != null
                          ? `${s.deep_minutes}m`
                          : "-"}
                      </td>
                      <td className="pr-6">
                        {s.rem_minutes != null ? `${s.rem_minutes}m` : "-"}
                      </td>
                      <td className="pr-6">
                        {s.efficiency != null ? `${s.efficiency}%` : "-"}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            ) : (
              <p className="text-sm text-(--color-fg-tertiary)">No sleep data</p>
            )}
          </div>

          {/* Heart Rate */}
          <div>
            <h3>Heart Rate</h3>
            {hrData && hrData.length > 0 ? (
              <table className="text-sm mt-2">
                <thead>
                  <tr className="text-(--color-fg-tertiary)">
                    <th className="pr-6 text-left">Date</th>
                    <th className="pr-6 text-left">Resting</th>
                    <th className="pr-6 text-left">Avg</th>
                    <th className="pr-6 text-left">Max</th>
                  </tr>
                </thead>
                <tbody>
                  {hrData.map((h) => (
                    <tr key={h.id}>
                      <td className="pr-6">
                        {new Date(h.date).toLocaleDateString()}
                      </td>
                      <td className="pr-6">
                        {h.resting_hr != null ? `${h.resting_hr} bpm` : "-"}
                      </td>
                      <td className="pr-6">
                        {h.avg_hr != null ? `${h.avg_hr} bpm` : "-"}
                      </td>
                      <td className="pr-6">
                        {h.max_hr != null ? `${h.max_hr} bpm` : "-"}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            ) : (
              <p className="text-sm text-(--color-fg-tertiary)">
                No heart rate data
              </p>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
