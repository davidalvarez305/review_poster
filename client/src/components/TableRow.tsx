import React, { useEffect, useState } from "react";
import { Td, Tr } from "@chakra-ui/react";
import TableButtons from "./TableButtons";

interface Props {
  items: any[];
  columns: string[];
  index: number;
  onClickEdit?: () => void;
  onClickDelete?: () => void;
}

const TableRow: React.FC<Props> = ({
  columns,
  index,
  items,
  onClickEdit,
  onClickDelete,
}) => {
  const [tableColumns, setTableColumns] = useState<string[]>(columns);
  useEffect(() => {
    setTableColumns(columns);
  }, [columns]);
  return (
    <Tr>
      {tableColumns.map((column, idx) => {
        return (
          <React.Fragment key={idx}>
            {column === "action" ? (
              <Td>
                <TableButtons
                  onClickEdit={onClickEdit}
                  onClickDelete={onClickDelete}
                />
              </Td>
            ) : (
              <Td>{items[index][column]}</Td>
            )}
          </React.Fragment>
        );
      })}
    </Tr>
  );
};

export default TableRow;
