import useLoginRequired from "./hooks/useLoginRequired";
import Layout from "./layout/Layout";

export const App = () => {
  useLoginRequired();
  return (
    <Layout>
      <div
        style={{
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          height: "80vh",
          width: "100%",
        }}
      >
        Welcome...
      </div>
    </Layout>
  );
};
