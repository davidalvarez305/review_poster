import { IconButton, IconButtonProps } from "@chakra-ui/react";
import React from "react";
import { FiSave } from "react-icons/fi";

export type Props = IconButtonProps & {
    "aria-label"?: string,
}

export const SaveButton: React.FC<Props> = ({ ...props }) => {
  return <IconButton colorScheme={'black'} variant={'outline'} {...props}  icon={<FiSave />} />;
};
