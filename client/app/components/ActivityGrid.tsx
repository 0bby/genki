import { useQuery } from "@tanstack/react-query";
import {
  getActivity,
  type getActivityArgs,
  type ActivityItem,
} from "api/api";
import Popup from "./Popup";
import { useState } from "react";
import { useTheme } from "~/hooks/useTheme";
import ActivityOptsSelector from "./ActivityOptsSelector";
import type { Theme } from "~/styles/themes.css";

const METRIC_COLORS: Record<string, string> = {
  workouts: "#e74c3c",
  steps: "#27ae60",
  sleep: "#3498db",
  active_minutes: "#e67e22",
};

const METRIC_LABELS: Record<string, string> = {
  workouts: "workouts",
  steps: "steps",
  sleep: "hours slept",
  active_minutes: "active min",
};

function getPrimaryColor(theme: Theme): string {
  const value = theme.primary;
  const rgbMatch = value.match(
    /^rgb\(\s*(\d{1,3})\s*,\s*(\d{1,3})\s*,\s*(\d{1,3})\s*\)$/
  );
  if (rgbMatch) {
    const [, r, g, b] = rgbMatch.map(Number);
    return "#" + [r, g, b].map((n) => n.toString(16).padStart(2, "0")).join("");
  }
  return value;
}

interface Props {
  step?: string;
  range?: number;
  month?: number;
  year?: number;
  configurable?: boolean;
}

export default function ActivityGrid({
  step = "day",
  range = 182,
  month = 0,
  year = 0,
  configurable = false,
}: Props) {
  const [stepState, setStep] = useState(step);
  const [rangeState, setRange] = useState(range);
  const [metric, setMetric] = useState("workouts");

  const { isPending, isError, data, error } = useQuery({
    queryKey: [
      "activity",
      {
        metric,
        step: stepState,
        range: rangeState,
        month: month,
        year: year,
      },
    ],
    queryFn: ({ queryKey }) => getActivity(queryKey[1] as getActivityArgs),
  });

  const color = METRIC_COLORS[metric] || "#e74c3c";

  if (isPending) {
    return (
      <div className="w-[350px]">
        <h3>Activity</h3>
        <p>Loading...</p>
      </div>
    );
  } else if (isError) {
    return (
      <div className="w-[350px]">
        <h3>Activity</h3>
        <p className="error">Error: {error.message}</p>
      </div>
    );
  }

  function LightenDarkenColor(hex: string, lum: number) {
    hex = String(hex).replace(/[^0-9a-f]/gi, "");
    if (hex.length < 6) {
      hex = hex[0] + hex[0] + hex[1] + hex[1] + hex[2] + hex[2];
    }
    lum = lum || 0;
    var rgb = "#", c, i;
    for (i = 0; i < 3; i++) {
      c = parseInt(hex.substring(i * 2, i * 2 + 2), 16);
      c = Math.round(Math.min(Math.max(0, c + c * lum), 255)).toString(16);
      rgb += ("00" + c).substring(c.length);
    }
    return rgb;
  }

  const getDarkenAmount = (v: number): number => {
    let t: number;
    switch (stepState) {
      case "day": t = metric === "steps" ? 10000 : 10; break;
      case "week": t = metric === "steps" ? 70000 : 20; break;
      case "month": t = metric === "steps" ? 300000 : 50; break;
      case "year": t = metric === "steps" ? 3000000 : 100; break;
      default: t = 10;
    }
    v = Math.min(v, t);
    return ((v - t) / t) * 0.8;
  };

  const CHUNK_SIZE = 26 * 7;
  const chunks = [];
  for (let i = 0; i < data.length; i += CHUNK_SIZE) {
    chunks.push(data.slice(i, i + CHUNK_SIZE));
  }

  const label = METRIC_LABELS[metric] || metric;

  return (
    <div className="flex flex-col items-start">
      <h3>Activity</h3>
      {configurable && (
        <>
          <div className="flex gap-2 mb-2 text-xs">
            {Object.keys(METRIC_COLORS).map((m) => (
              <button
                key={m}
                className={`px-2 py-0.5 rounded transition ${
                  m === metric
                    ? "font-bold border border-current"
                    : "color-fg-secondary hover:color-fg"
                }`}
                style={m === metric ? { color: METRIC_COLORS[m] } : undefined}
                onClick={() => setMetric(m)}
              >
                {m.replace("_", " ")}
              </button>
            ))}
          </div>
          <ActivityOptsSelector
            rangeSetter={setRange}
            currentRange={rangeState}
            stepSetter={setStep}
            currentStep={stepState}
          />
        </>
      )}

      {chunks.map((chunk, index) => (
        <div
          key={index}
          className="w-auto grid grid-flow-col grid-rows-7 gap-[3px] md:gap-[5px] mb-4"
        >
          {chunk.map((item) => (
            <div
              key={item.start}
              className="w-[10px] sm:w-[12px] h-[10px] sm:h-[12px]"
            >
              <Popup
                position="top"
                space={12}
                extraClasses="left-2"
                inner={`${new Date(item.start).toLocaleDateString()} — ${
                  metric === "sleep"
                    ? (item.value / 60).toFixed(1) + "h"
                    : metric === "steps"
                    ? item.value.toLocaleString()
                    : item.value
                } ${label}`}
              >
                <div
                  style={{
                    display: "inline-block",
                    background:
                      item.value > 0
                        ? LightenDarkenColor(color, getDarkenAmount(item.value))
                        : "var(--color-bg-secondary)",
                  }}
                  className={`w-[10px] sm:w-[12px] h-[10px] sm:h-[12px] rounded-[2px] md:rounded-[3px] ${
                    item.value > 0
                      ? ""
                      : "border-[0.5px] border-(--color-bg-tertiary)"
                  }`}
                ></div>
              </Popup>
            </div>
          ))}
        </div>
      ))}
    </div>
  );
}
