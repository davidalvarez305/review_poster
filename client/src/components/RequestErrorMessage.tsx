import { Box } from "@chakra-ui/layout";
import React from "react";

interface Props {
  message?: string;
}

const RequestErrorMessage: React.FC<Props> = ({ message }) => {
  return (
    <React.Fragment>
      {message && message.length > 0 && (
        <Box
          sx={{
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
          }}
        >
          <p
            style={{
              fontFamily: "Georgia",
              color: "red",
              fontSize: 18,
            }}
          >
            {message}
          </p>
        </Box>
      )}
    </React.Fragment>
  );
};

export default RequestErrorMessage;
