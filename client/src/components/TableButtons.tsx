import { Box, IconButton } from "@chakra-ui/react";
import React from "react";
import { EditIcon, DeleteIcon } from "@chakra-ui/icons";
type Props = {
  onClickEdit?: () => void;
  onClickDelete?: () => void;
};

const TableButtons: React.FC<Props> = ({ onClickEdit, onClickDelete }) => {
  return (
    <Box
      sx={{
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        gap: 2,
        my: 2,
      }}
    >
      <IconButton aria-label="edit" icon={<EditIcon />} onClick={onClickEdit} />
      <IconButton
        aria-label="delete"
        icon={<DeleteIcon />}
        onClick={onClickDelete}
      />
    </Box>
  );
};

export default TableButtons;
