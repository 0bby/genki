import Rewind from "~/components/rewind/Rewind";
import type { RecapStats } from "api/api";
import { useState } from "react";
import type { LoaderFunctionArgs } from "react-router";
import { useLoaderData, useNavigate } from "react-router";
import { getRewindParams } from "~/utils/utils";
import { ChevronLeft, ChevronRight } from "lucide-react";

const months = [
  "Full Year",
  "January", "February", "March", "April", "May", "June",
  "July", "August", "September", "October", "November", "December",
];

export async function clientLoader({ request }: LoaderFunctionArgs) {
  const url = new URL(request.url);
  const year = parseInt(url.searchParams.get("year") || getRewindParams().year.toString());
  const month = parseInt(url.searchParams.get("month") || getRewindParams().month.toString());

  const res = await fetch(`/apis/web/v1/summary?year=${year}&month=${month}`);
  if (!res.ok) {
    throw new Response("Failed to load summary", { status: 500 });
  }

  const stats: RecapStats = await res.json();
  stats.title = `Your ${month === 0 ? "" : months[month] + " "}${year} Recap`;
  return { stats };
}

export default function RecapPage() {
  const currentParams = new URLSearchParams(location.search);
  let year = parseInt(currentParams.get("year") || getRewindParams().year.toString());
  let month = parseInt(currentParams.get("month") || getRewindParams().month.toString());
  const navigate = useNavigate();
  const { stats } = useLoaderData<{ stats: RecapStats }>();

  const updateParams = (params: Record<string, string | null>) => {
    const nextParams = new URLSearchParams(location.search);
    for (const key in params) {
      const val = params[key];
      if (val !== null) nextParams.set(key, val);
    }
    navigate(`/recap?${nextParams.toString()}`, { replace: false });
  };

  const navigateMonth = (direction: "prev" | "next") => {
    if (direction === "next") { month = month === 12 ? 0 : month + 1; }
    else { month = month === 0 ? 12 : month - 1; }
    updateParams({ year: year.toString(), month: month.toString() });
  };

  const navigateYear = (direction: "prev" | "next") => {
    year += direction === "next" ? 1 : -1;
    updateParams({ year: year.toString(), month: month.toString() });
  };

  const pgTitle = `${stats.title} - Genki`;

  return (
    <div className="w-full min-h-screen">
      <div className="flex flex-col items-start sm:items-center gap-4">
        <title>{pgTitle}</title>
        <div className="flex flex-col lg:flex-row items-start lg:mt-15 mt-5 gap-10 w-19/20 px-5 md:px-20">
          <div className="flex flex-col items-start gap-4 py-8">
            <div className="flex items-center gap-6 justify-around">
              <button onClick={() => navigateMonth("prev")} className="p-2">
                <ChevronLeft size={20} />
              </button>
              <p className="font-medium text-xl text-center w-30">{months[month]}</p>
              <button onClick={() => navigateMonth("next")} className="p-2">
                <ChevronRight size={20} />
              </button>
            </div>
            <div className="flex items-center gap-6 justify-around">
              <button onClick={() => navigateYear("prev")} className="p-2">
                <ChevronLeft size={20} />
              </button>
              <p className="font-medium text-xl text-center w-30">{year}</p>
              <button onClick={() => navigateYear("next")} className="p-2">
                <ChevronRight size={20} />
              </button>
            </div>
          </div>
          {stats && <Rewind stats={stats} />}
        </div>
      </div>
    </div>
  );
}
