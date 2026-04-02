import { useQuery } from "@tanstack/react-query";
import { getTopMuscles, type getItemsArgs } from "api/api";
import { Link } from "react-router";
import TopItemList from "./TopItemList";

interface Props {
  limit: number;
  period: string;
}

export default function TopMuscles(props: Props) {
  const { isPending, isError, data, error } = useQuery({
    queryKey: [
      "top-muscles",
      { limit: props.limit, period: props.period, page: 0 },
    ],
    queryFn: ({ queryKey }) => getTopMuscles(queryKey[1] as getItemsArgs),
  });

  const header = "Top muscles";

  if (isPending) {
    return (
      <div className="w-[300px]">
        <h3>{header}</h3>
        <p>Loading...</p>
      </div>
    );
  } else if (isError) {
    return (
      <div className="w-[300px]">
        <h3>{header}</h3>
        <p className="error">Error: {error.message}</p>
      </div>
    );
  }

  return (
    <div>
      <h3 className="hover:underline">
        <Link to={`/chart/top-muscles?period=${props.period}`}>{header}</Link>
      </h3>
      <div className="max-w-[300px]">
        <TopItemList type="muscle" data={data} />
        {data.items.length < 1 ? "Nothing to show" : ""}
      </div>
    </div>
  );
}
