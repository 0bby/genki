import { useState } from "react";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { getSyncStatus, triggerSync, initFitbitOAuth, type SyncStatus } from "api/api";
import { AsyncButton } from "../AsyncButton";

export default function DataSourcesModal() {
  const queryClient = useQueryClient();
  const [wgerSyncing, setWgerSyncing] = useState(false);
  const [fitbitSyncing, setFitbitSyncing] = useState(false);
  const [fitbitConnecting, setFitbitConnecting] = useState(false);
  const [message, setMessage] = useState("");

  const { data: syncStatus } = useQuery({
    queryKey: ["sync-status"],
    queryFn: getSyncStatus,
    refetchInterval: 10000,
  });

  const handleWgerSync = async () => {
    setWgerSyncing(true);
    setMessage("");
    try {
      const r = await triggerSync("wger");
      if (r.ok) {
        setMessage("wger sync triggered");
        setTimeout(() => queryClient.invalidateQueries({ queryKey: ["sync-status"] }), 3000);
      } else {
        setMessage("Failed to trigger wger sync");
      }
    } catch {
      setMessage("Failed to trigger wger sync");
    }
    setWgerSyncing(false);
  };

  const handleFitbitSync = async () => {
    setFitbitSyncing(true);
    setMessage("");
    try {
      const r = await triggerSync("fitbit");
      if (r.ok) {
        setMessage("Fitbit sync triggered");
        setTimeout(() => queryClient.invalidateQueries({ queryKey: ["sync-status"] }), 3000);
      } else {
        setMessage("Failed to trigger Fitbit sync");
      }
    } catch {
      setMessage("Failed to trigger Fitbit sync");
    }
    setFitbitSyncing(false);
  };

  const handleFitbitConnect = async () => {
    setFitbitConnecting(true);
    setMessage("");
    try {
      const data = await initFitbitOAuth();
      window.location.href = data.url;
    } catch {
      setMessage("Failed to start Fitbit connection");
      setFitbitConnecting(false);
    }
  };

  const hasFitbitToken = syncStatus?.sources?.some(
    (s) => s.source === "fitbit"
  );

  const formatDate = (dateStr: string) => {
    const d = new Date(dateStr);
    return d.toLocaleString();
  };

  return (
    <div className="flex flex-col gap-6">
      <h3>Data Sources</h3>

      <div className="flex flex-col gap-4">
        <div className="border border-(--color-bg-tertiary) rounded-lg p-4">
          <div className="flex items-center justify-between mb-2">
            <h4 className="font-medium">wger</h4>
            <AsyncButton loading={wgerSyncing} onClick={handleWgerSync}>
              Sync Now
            </AsyncButton>
          </div>
          <p className="text-sm color-fg-secondary mb-2">
            Syncs exercises, workouts, and sets from your wger instance.
            Runs automatically every 5 minutes when configured.
          </p>
          {renderCursors(syncStatus, "wger", formatDate)}
        </div>

        <div className="border border-(--color-bg-tertiary) rounded-lg p-4">
          <div className="flex items-center justify-between mb-2">
            <h4 className="font-medium">Fitbit</h4>
            <div className="flex gap-2">
              {hasFitbitToken ? (
                <AsyncButton loading={fitbitSyncing} onClick={handleFitbitSync}>
                  Sync Now
                </AsyncButton>
              ) : (
                <AsyncButton loading={fitbitConnecting} onClick={handleFitbitConnect}>
                  Connect Fitbit
                </AsyncButton>
              )}
            </div>
          </div>
          <p className="text-sm color-fg-secondary mb-2">
            {hasFitbitToken
              ? "Connected. Syncs steps, sleep, heart rate, and activity."
              : "Connect your Fitbit account to sync steps, sleep, heart rate, and activity data."}
          </p>
          {renderCursors(syncStatus, "fitbit", formatDate)}
        </div>
      </div>

      {message && (
        <p className="text-sm color-fg-secondary">{message}</p>
      )}
    </div>
  );
}

function renderCursors(
  syncStatus: SyncStatus | undefined,
  source: string,
  formatDate: (s: string) => string
) {
  const cursors = syncStatus?.sources?.filter((s) => s.source === source);
  if (!cursors || cursors.length === 0) return null;

  return (
    <div className="text-xs color-fg-tertiary mt-2">
      {cursors.map((c) => (
        <div key={`${c.source}-${c.resource}`}>
          {c.resource}: last synced {formatDate(c.last_synced_at)}
        </div>
      ))}
    </div>
  );
}
