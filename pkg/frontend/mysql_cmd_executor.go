package frontend

import (
	"fmt"
	"matrixone/pkg/defines"
	"matrixone/pkg/sql/compile"
	"strings"
	"time"

	"matrixone/pkg/container/batch"
	"matrixone/pkg/container/types"
	"matrixone/pkg/sql/tree"
	"matrixone/pkg/vm/process"
)

type MysqlCmdExecutor struct {
	CmdExecutorImpl

	//the count of sql has been processed
	sqlCount uint64
}

//get new process id
func (mce *MysqlCmdExecutor) getNextProcessId() string {
	/*
		temporary method:
		routineId + sqlCount
	*/
	routineId := mce.routine.getConnID()
	return fmt.Sprintf("%d%d", routineId, mce.sqlCount)
}

func (mce *MysqlCmdExecutor) addSqlCount(a uint64) {
	mce.sqlCount += a
}

/*
extract the data from the pipeline.
obj: routine obj
TODO:Add error
Warning: The pipeline is the multi-thread environment. The getDataFromPipeline will
	access the shared data. Be careful when it writes the shared data.
*/
func getDataFromPipeline(obj interface{}, bat *batch.Batch) error {
	rt := obj.(*Routine)
	ses := rt.GetSession()

	fmt.Println("hello------")
	{
		fmt.Printf("bat: %v\n", bat)
	}

	var rowGroupSize = ses.Pu.SV.GetCountOfRowsPerSendingToClient()
	rowGroupSize = MaxInt64(rowGroupSize, 1)

	var choose = !ses.Pu.SV.GetSendRow()
	if choose {
		goID := GetRoutineId()

		fmt.Printf("goid %d \n", goID)

		proto := rt.GetClientProtocol().(*MysqlProtocol)

		//Create a new temporary resultset per pipeline thread.
		mrs := &MysqlResultSet{}
		//Warning: Don't change Columns in this.
		//Reference the shared Columns of the session among multi-thread.
		mrs.Columns = ses.Mrs.Columns
		mrs.Name2Index = ses.Mrs.Name2Index

		//group row
		mrs.Data = make([][]interface{}, rowGroupSize)
		for i := int64(0); i < rowGroupSize; i++ {
			mrs.Data[i] = make([]interface{}, len(bat.Vecs))
		}

		if n := len(bat.Sels); n == 0 {
			n = bat.Vecs[0].Length()
			groupCnt := int64(n)/rowGroupSize + 1
			for g := int64(0); g < groupCnt; g++ { //group id
				begin := g * rowGroupSize
				end := MinInt64((g+1)*rowGroupSize, int64(n))

				r := uint64(0)
				for j := begin; j < end; j++ { //row index
					row := mrs.Data[r]
					r++
					for i, vec := range bat.Vecs { //col index
						switch vec.Typ.Oid { //get col
						case types.T_int8:
							if !vec.Nsp.Any() { //all data in this column are not null
								vs := vec.Col.([]int8)
								row[i] = vs[j]
							} else {
								if vec.Nsp.Contains(uint64(j)) { //is null
									row[i] = nil
								} else {
									vs := vec.Col.([]int8)
									row[i] = vs[j]
								}
							}
						case types.T_uint8:
							if !vec.Nsp.Any() { //all data in this column are not null
								vs := vec.Col.([]uint8)
								row[i] = vs[j]
							} else {
								if vec.Nsp.Contains(uint64(j)) { //is null
									row[i] = nil
								} else {
									vs := vec.Col.([]uint8)
									row[i] = vs[j]
								}
							}
						case types.T_int16:
							if !vec.Nsp.Any() { //all data in this column are not null
								vs := vec.Col.([]int16)
								row[i] = vs[j]
							} else {
								if vec.Nsp.Contains(uint64(j)) { //is null
									row[i] = nil
								} else {
									vs := vec.Col.([]int16)
									row[i] = vs[j]
								}
							}
						case types.T_uint16:
							if !vec.Nsp.Any() { //all data in this column are not null
								vs := vec.Col.([]uint16)
								row[i] = vs[j]
							} else {
								if vec.Nsp.Contains(uint64(j)) { //is null
									row[i] = nil
								} else {
									vs := vec.Col.([]uint16)
									row[i] = vs[j]
								}
							}
						case types.T_int32:
							if !vec.Nsp.Any() { //all data in this column are not null
								vs := vec.Col.([]int32)
								row[i] = vs[j]
							} else {
								if vec.Nsp.Contains(uint64(j)) { //is null
									row[i] = nil
								} else {
									vs := vec.Col.([]int32)
									row[i] = vs[j]
								}
							}
						case types.T_uint32:
							if !vec.Nsp.Any() { //all data in this column are not null
								vs := vec.Col.([]uint32)
								row[i] = vs[j]
							} else {
								if vec.Nsp.Contains(uint64(j)) { //is null
									row[i] = nil
								} else {
									vs := vec.Col.([]uint32)
									row[i] = vs[j]
								}
							}
						case types.T_int64:
							if !vec.Nsp.Any() { //all data in this column are not null
								vs := vec.Col.([]int64)
								row[i] = vs[j]
							} else {
								if vec.Nsp.Contains(uint64(j)) { //is null
									row[i] = nil
								} else {
									vs := vec.Col.([]int64)
									row[i] = vs[j]
								}
							}
						case types.T_uint64:
							if !vec.Nsp.Any() { //all data in this column are not null
								vs := vec.Col.([]uint64)
								row[i] = vs[j]
							} else {
								if vec.Nsp.Contains(uint64(j)) { //is null
									row[i] = nil
								} else {
									vs := vec.Col.([]uint64)
									row[i] = vs[j]
								}
							}
						case types.T_float32:
							if !vec.Nsp.Any() { //all data in this column are not null
								vs := vec.Col.([]float32)
								row[i] = vs[j]
							} else {
								if vec.Nsp.Contains(uint64(j)) { //is null
									row[i] = nil
								} else {
									vs := vec.Col.([]float32)
									row[i] = vs[j]
								}
							}
						case types.T_float64:
							if !vec.Nsp.Any() { //all data in this column are not null
								vs := vec.Col.([]float64)
								row[i] = vs[j]
							} else {
								if vec.Nsp.Contains(uint64(j)) { //is null
									row[i] = nil
								} else {
									vs := vec.Col.([]float64)
									row[i] = vs[j]
								}
							}
						case types.T_char:
							if !vec.Nsp.Any() { //all data in this column are not null
								vs := vec.Col.(*types.Bytes)
								row[i] = vs.Get(j)
							} else {
								if vec.Nsp.Contains(uint64(j)) { //is null
									row[i] = nil
								} else {
									vs := vec.Col.(*types.Bytes)
									row[i] = vs.Get(j)
								}
							}
						case types.T_varchar:
							if !vec.Nsp.Any() { //all data in this column are not null
								vs := vec.Col.(*types.Bytes)
								row[i] = vs.Get(j)
							} else {
								if vec.Nsp.Contains(uint64(j)) { //is null
									row[i] = nil
								} else {
									vs := vec.Col.(*types.Bytes)
									row[i] = vs.Get(j)
								}
							}
						default:
							fmt.Printf("getDataFromPipeline : unsupported type %d \n", vec.Typ.Oid)
							return fmt.Errorf("getDataFromPipeline : unsupported type %d \n", vec.Typ.Oid)
						}
					}
				}

				fmt.Printf("row group -+> %v \n", mrs.Data[:r])

				//send group of row
				if err := proto.SendResultSetTextBatchRow(mrs, r); err != nil {
					//return err
					fmt.Printf("getDataFromPipeline error %v \n", err)
					return err
				}
			}

		} else {
			n = bat.Vecs[0].Length()
			groupCnt := int64(n)/rowGroupSize + 1
			for g := int64(0); g < groupCnt; g++ { //group id
				begin := g * rowGroupSize
				end := MinInt64((g+1)*rowGroupSize, int64(n))

				r := uint64(0)
				for j := begin; j < end; j++ { //row index
					row := mrs.Data[r]
					r++
					for i, vec := range bat.Vecs { //col index
						switch vec.Typ.Oid { //get col
						case types.T_int8:
							if !vec.Nsp.Any() { //all data in this column are not null
								vs := vec.Col.([]int8)
								row[i] = vs[bat.Sels[j]]
							} else {
								if vec.Nsp.Contains(uint64(j)) { //is null
									row[i] = nil
								} else {
									vs := vec.Col.([]int8)
									row[i] = vs[bat.Sels[j]]
								}
							}
						case types.T_uint8:
							if !vec.Nsp.Any() { //all data in this column are not null
								vs := vec.Col.([]uint8)
								row[i] = vs[bat.Sels[j]]
							} else {
								if vec.Nsp.Contains(uint64(j)) { //is null
									row[i] = nil
								} else {
									vs := vec.Col.([]uint8)
									row[i] = vs[bat.Sels[j]]
								}
							}
						case types.T_int16:
							if !vec.Nsp.Any() { //all data in this column are not null
								vs := vec.Col.([]int16)
								row[i] = vs[bat.Sels[j]]
							} else {
								if vec.Nsp.Contains(uint64(j)) { //is null
									row[i] = nil
								} else {
									vs := vec.Col.([]int16)
									row[i] = vs[bat.Sels[j]]
								}
							}
						case types.T_uint16:
							if !vec.Nsp.Any() { //all data in this column are not null
								vs := vec.Col.([]uint16)
								row[i] = vs[bat.Sels[j]]
							} else {
								if vec.Nsp.Contains(uint64(j)) { //is null
									row[i] = nil
								} else {
									vs := vec.Col.([]uint16)
									row[i] = vs[bat.Sels[j]]
								}
							}
						case types.T_int32:
							if !vec.Nsp.Any() { //all data in this column are not null
								vs := vec.Col.([]int32)
								row[i] = vs[bat.Sels[j]]
							} else {
								if vec.Nsp.Contains(uint64(j)) { //is null
									row[i] = nil
								} else {
									vs := vec.Col.([]int32)
									row[i] = vs[bat.Sels[j]]
								}
							}
						case types.T_uint32:
							if !vec.Nsp.Any() { //all data in this column are not null
								vs := vec.Col.([]uint32)
								row[i] = vs[bat.Sels[j]]
							} else {
								if vec.Nsp.Contains(uint64(j)) { //is null
									row[i] = nil
								} else {
									vs := vec.Col.([]uint32)
									row[i] = vs[bat.Sels[j]]
								}
							}
						case types.T_int64:
							if !vec.Nsp.Any() { //all data in this column are not null
								vs := vec.Col.([]int64)
								row[i] = vs[bat.Sels[j]]
							} else {
								if vec.Nsp.Contains(uint64(j)) { //is null
									row[i] = nil
								} else {
									vs := vec.Col.([]int64)
									row[i] = vs[bat.Sels[j]]
								}
							}
						case types.T_uint64:
							if !vec.Nsp.Any() { //all data in this column are not null
								vs := vec.Col.([]uint64)
								row[i] = vs[bat.Sels[j]]
							} else {
								if vec.Nsp.Contains(uint64(j)) { //is null
									row[i] = nil
								} else {
									vs := vec.Col.([]uint64)
									row[i] = vs[bat.Sels[j]]
								}
							}
						case types.T_float32:
							if !vec.Nsp.Any() { //all data in this column are not null
								vs := vec.Col.([]float32)
								row[i] = vs[bat.Sels[j]]
							} else {
								if vec.Nsp.Contains(uint64(j)) { //is null
									row[i] = nil
								} else {
									vs := vec.Col.([]float32)
									row[i] = vs[bat.Sels[j]]
								}
							}
						case types.T_float64:
							if !vec.Nsp.Any() { //all data in this column are not null
								vs := vec.Col.([]float64)
								row[i] = vs[bat.Sels[j]]
							} else {
								if vec.Nsp.Contains(uint64(j)) { //is null
									row[i] = nil
								} else {
									vs := vec.Col.([]float64)
									row[i] = vs[bat.Sels[j]]
								}
							}
						case types.T_char:
							if !vec.Nsp.Any() { //all data in this column are not null
								vs := vec.Col.(*types.Bytes)
								row[i] = vs.Get(bat.Sels[j])
							} else {
								if vec.Nsp.Contains(uint64(j)) { //is null
									row[i] = nil
								} else {
									vs := vec.Col.(*types.Bytes)
									row[i] = vs.Get(bat.Sels[j])
								}
							}
						case types.T_varchar:
							if !vec.Nsp.Any() { //all data in this column are not null
								vs := vec.Col.(*types.Bytes)
								row[i] = vs.Get(bat.Sels[j])
							} else {
								if vec.Nsp.Contains(uint64(j)) { //is null
									row[i] = nil
								} else {
									vs := vec.Col.(*types.Bytes)
									row[i] = vs.Get(bat.Sels[j])
								}
							}
						default:
							fmt.Printf("getDataFromPipeline : unsupported type %d \n", vec.Typ.Oid)
							return fmt.Errorf("getDataFromPipeline : unsupported type %d \n", vec.Typ.Oid)
						}
					}
				}

				fmt.Printf("row group -*> %v \n", mrs.Data[:r])

				//send row
				if err := proto.SendResultSetTextBatchRow(mrs, r); err != nil {
					//return err
					fmt.Printf("getDataFromPipeline error %v \n", err)
					return err
				}
			}
		}
	} else {

		if n := len(bat.Sels); n == 0 {
			n = bat.Vecs[0].Length()
			rows := make([][]interface{}, n)
			for i := 0; i < n; i++ {
				rows[i] = make([]interface{}, len(bat.Vecs))
			}
			for i, vec := range bat.Vecs { //column index
				switch vec.Typ.Oid {
				case types.T_int8:
					if !vec.Nsp.Any() { //all data in this column are not null
						vs := vec.Col.([]int8)
						for j := 0; j < n; j++ {
							rows[j][i] = vs[j]
						}
					} else {
						for j := 0; j < n; j++ {
							if vec.Nsp.Contains(uint64(j)) { //is null
								rows[j][i] = nil
							} else {
								vs := vec.Col.([]int8)
								rows[j][i] = vs[j]
							}
						}
					}
				case types.T_uint8:
					if !vec.Nsp.Any() { //all data in this column are not null
						vs := vec.Col.([]uint8)
						for j := 0; j < n; j++ {
							rows[j][i] = vs[j]
						}
					} else {
						for j := 0; j < n; j++ {
							if vec.Nsp.Contains(uint64(j)) { //is null
								rows[j][i] = nil
							} else {
								vs := vec.Col.([]uint8)
								rows[j][i] = vs[j]
							}
						}
					}
				case types.T_int16:
					if !vec.Nsp.Any() { //all data in this column are not null
						vs := vec.Col.([]int16)
						for j := 0; j < n; j++ {
							rows[j][i] = vs[j]
						}
					} else {
						for j := 0; j < n; j++ {
							if vec.Nsp.Contains(uint64(j)) { //is null
								rows[j][i] = nil
							} else {
								vs := vec.Col.([]int16)
								rows[j][i] = vs[j]
							}
						}
					}
				case types.T_uint16:
					if !vec.Nsp.Any() { //all data in this column are not null
						vs := vec.Col.([]uint16)
						for j := 0; j < n; j++ {
							rows[j][i] = vs[j]
						}
					} else {
						for j := 0; j < n; j++ {
							if vec.Nsp.Contains(uint64(j)) { //is null
								rows[j][i] = nil
							} else {
								vs := vec.Col.([]uint16)
								rows[j][i] = vs[j]
							}
						}
					}
				case types.T_int32:
					if !vec.Nsp.Any() { //all data in this column are not null
						vs := vec.Col.([]int32)
						for j := 0; j < n; j++ {
							rows[j][i] = vs[j]
						}
					} else {
						for j := 0; j < n; j++ {
							if vec.Nsp.Contains(uint64(j)) { //is null
								rows[j][i] = nil
							} else {
								vs := vec.Col.([]int32)
								rows[j][i] = vs[j]
							}
						}
					}
				case types.T_uint32:
					if !vec.Nsp.Any() { //all data in this column are not null
						vs := vec.Col.([]uint32)
						for j := 0; j < n; j++ {
							rows[j][i] = vs[j]
						}
					} else {
						for j := 0; j < n; j++ {
							if vec.Nsp.Contains(uint64(j)) { //is null
								rows[j][i] = nil
							} else {
								vs := vec.Col.([]uint32)
								rows[j][i] = vs[j]
							}
						}
					}
				case types.T_int64:
					if !vec.Nsp.Any() { //all data in this column are not null
						vs := vec.Col.([]int64)
						for j := 0; j < n; j++ {
							rows[j][i] = vs[j]
						}
					} else {
						for j := 0; j < n; j++ {
							if vec.Nsp.Contains(uint64(j)) { //is null
								rows[j][i] = nil
							} else {
								vs := vec.Col.([]int64)
								rows[j][i] = vs[j]
							}
						}
					}
				case types.T_uint64:
					if !vec.Nsp.Any() { //all data in this column are not null
						vs := vec.Col.([]uint64)
						for j := 0; j < n; j++ {
							rows[j][i] = vs[j]
						}
					} else {
						for j := 0; j < n; j++ {
							if vec.Nsp.Contains(uint64(j)) { //is null
								rows[j][i] = nil
							} else {
								vs := vec.Col.([]uint64)
								rows[j][i] = vs[j]
							}
						}
					}
				case types.T_float32:
					if !vec.Nsp.Any() { //all data in this column are not null
						vs := vec.Col.([]float32)
						for j := 0; j < n; j++ {
							rows[j][i] = vs[j]
						}
					} else {
						for j := 0; j < n; j++ {
							if vec.Nsp.Contains(uint64(j)) { //is null
								rows[j][i] = nil
							} else {
								vs := vec.Col.([]float32)
								rows[j][i] = vs[j]
							}
						}
					}
				case types.T_float64:
					if !vec.Nsp.Any() { //all data in this column are not null
						vs := vec.Col.([]float64)
						for j := 0; j < n; j++ {
							rows[j][i] = vs[j]
						}
					} else {
						for j := 0; j < n; j++ {
							if vec.Nsp.Contains(uint64(j)) { //is null
								rows[j][i] = nil
							} else {
								vs := vec.Col.([]float64)
								rows[j][i] = vs[j]
							}
						}
					}
				case types.T_char:
					if !vec.Nsp.Any() { //all data in this column are not null
						vs := vec.Col.(*types.Bytes)
						for j := 0; j < n; j++ {
							rows[j][i] = vs.Get(int64(j))
						}
					} else {
						for j := 0; j < n; j++ {
							if vec.Nsp.Contains(uint64(j)) { //is null
								rows[j][i] = nil
							} else {
								vs := vec.Col.(*types.Bytes)
								rows[j][i] = vs.Get(int64(j))
							}
						}
					}
				case types.T_varchar:
					if !vec.Nsp.Any() { //all data in this column are not null
						vs := vec.Col.(*types.Bytes)
						for j := 0; j < n; j++ {
							rows[j][i] = vs.Get(int64(j))
						}
					} else {
						for j := 0; j < n; j++ {
							if vec.Nsp.Contains(uint64(j)) { //is null
								rows[j][i] = nil
							} else {
								vs := vec.Col.(*types.Bytes)
								rows[j][i] = vs.Get(int64(j))
							}
						}
					}
				default:
					fmt.Printf("FillResult else1: unsupported type %d \n", vec.Typ.Oid)
					return fmt.Errorf("FillResult else1: unsupported type %d \n", vec.Typ.Oid)
				}
			}
			ses.Mrs.Data = rows
		} else {
			rows := make([][]interface{}, n)
			for i := 0; i < n; i++ {
				rows[i] = make([]interface{}, len(bat.Vecs))
			}
			for i, vec := range bat.Vecs {
				switch vec.Typ.Oid {
				case types.T_int8:
					if !vec.Nsp.Any() { //all data in this column are not null
						vs := vec.Col.([]int8)
						for j := 0; j < n; j++ {
							rows[j][i] = vs[bat.Sels[j]]
						}
					} else {
						for j := 0; j < n; j++ {
							if vec.Nsp.Contains(uint64(j)) { //is null
								rows[j][i] = nil
							} else {
								vs := vec.Col.([]int8)
								rows[j][i] = vs[bat.Sels[j]]
							}
						}
					}
				case types.T_uint8:
					if !vec.Nsp.Any() { //all data in this column are not null
						vs := vec.Col.([]uint8)
						for j := 0; j < n; j++ {
							rows[j][i] = vs[bat.Sels[j]]
						}
					} else {
						for j := 0; j < n; j++ {
							if vec.Nsp.Contains(uint64(j)) { //is null
								rows[j][i] = nil
							} else {
								vs := vec.Col.([]uint8)
								rows[j][i] = vs[bat.Sels[j]]
							}
						}
					}
				case types.T_int16:
					if !vec.Nsp.Any() { //all data in this column are not null
						vs := vec.Col.([]int16)
						for j := 0; j < n; j++ {
							rows[j][i] = vs[bat.Sels[j]]
						}
					} else {
						for j := 0; j < n; j++ {
							if vec.Nsp.Contains(uint64(j)) { //is null
								rows[j][i] = nil
							} else {
								vs := vec.Col.([]int16)
								rows[j][i] = vs[bat.Sels[j]]
							}
						}
					}
				case types.T_uint16:
					if !vec.Nsp.Any() { //all data in this column are not null
						vs := vec.Col.([]uint16)
						for j := 0; j < n; j++ {
							rows[j][i] = vs[bat.Sels[j]]
						}
					} else {
						for j := 0; j < n; j++ {
							if vec.Nsp.Contains(uint64(j)) { //is null
								rows[j][i] = nil
							} else {
								vs := vec.Col.([]uint16)
								rows[j][i] = vs[bat.Sels[j]]
							}
						}
					}
				case types.T_int32:
					if !vec.Nsp.Any() { //all data in this column are not null
						vs := vec.Col.([]int32)
						for j := 0; j < n; j++ {
							rows[j][i] = vs[bat.Sels[j]]
						}
					} else {
						for j := 0; j < n; j++ {
							if vec.Nsp.Contains(uint64(j)) { //is null
								rows[j][i] = nil
							} else {
								vs := vec.Col.([]int32)
								rows[j][i] = vs[bat.Sels[j]]
							}
						}
					}
				case types.T_uint32:
					if !vec.Nsp.Any() { //all data in this column are not null
						vs := vec.Col.([]uint32)
						for j := 0; j < n; j++ {
							rows[j][i] = vs[bat.Sels[j]]
						}
					} else {
						for j := 0; j < n; j++ {
							if vec.Nsp.Contains(uint64(j)) { //is null
								rows[j][i] = nil
							} else {
								vs := vec.Col.([]uint32)
								rows[j][i] = vs[bat.Sels[j]]
							}
						}
					}
				case types.T_int64:
					if !vec.Nsp.Any() { //all data in this column are not null
						vs := vec.Col.([]int64)
						for j := 0; j < n; j++ {
							rows[j][i] = vs[bat.Sels[j]]
						}
					} else {
						for j := 0; j < n; j++ {
							if vec.Nsp.Contains(uint64(j)) { //is null
								rows[j][i] = nil
							} else {
								vs := vec.Col.([]int64)
								rows[j][i] = vs[bat.Sels[j]]
							}
						}
					}
				case types.T_uint64:
					if !vec.Nsp.Any() { //all data in this column are not null
						vs := vec.Col.([]uint64)
						for j := 0; j < n; j++ {
							rows[j][i] = vs[bat.Sels[j]]
						}
					} else {
						for j := 0; j < n; j++ {
							if vec.Nsp.Contains(uint64(j)) { //is null
								rows[j][i] = nil
							} else {
								vs := vec.Col.([]uint64)
								rows[j][i] = vs[bat.Sels[j]]
							}
						}
					}
				case types.T_float32:
					if !vec.Nsp.Any() { //all data in this column are not null
						vs := vec.Col.([]float32)
						for j := 0; j < n; j++ {
							rows[j][i] = vs[bat.Sels[j]]
						}
					} else {
						for j := 0; j < n; j++ {
							if vec.Nsp.Contains(uint64(j)) { //is null
								rows[j][i] = nil
							} else {
								vs := vec.Col.([]float32)
								rows[j][i] = vs[bat.Sels[j]]
							}
						}
					}
				case types.T_float64:
					if !vec.Nsp.Any() { //all data in this column are not null
						vs := vec.Col.([]float64)
						for j := 0; j < n; j++ {
							rows[j][i] = vs[bat.Sels[j]]
						}
					} else {
						for j := 0; j < n; j++ {
							if vec.Nsp.Contains(uint64(j)) { //is null
								rows[j][i] = nil
							} else {
								vs := vec.Col.([]float64)
								rows[j][i] = vs[bat.Sels[j]]
							}
						}
					}
				case types.T_char:
					if !vec.Nsp.Any() { //all data in this column are not null
						vs := vec.Col.(*types.Bytes)
						for j := 0; j < n; j++ {
							rows[j][i] = vs.Get(bat.Sels[j])
						}
					} else {
						for j := 0; j < n; j++ {
							if vec.Nsp.Contains(uint64(j)) { //is null
								rows[j][i] = nil
							} else {
								vs := vec.Col.(*types.Bytes)
								rows[j][i] = vs.Get(bat.Sels[j])
							}
						}
					}
				case types.T_varchar:
					if !vec.Nsp.Any() { //all data in this column are not null
						vs := vec.Col.(*types.Bytes)
						for j := 0; j < n; j++ {
							rows[j][i] = vs.Get(bat.Sels[j])
						}
					} else {
						for j := 0; j < n; j++ {
							if vec.Nsp.Contains(uint64(j)) { //is null
								rows[j][i] = nil
							} else {
								vs := vec.Col.(*types.Bytes)
								rows[j][i] = vs.Get(bat.Sels[j])
							}
						}
					}
				default:
					fmt.Printf("FillResult else2: unsupported type %d \n", vec.Typ.Oid)
					return fmt.Errorf("FillResult else2: unsupported type %d \n", vec.Typ.Oid)
				}
			}
		}
	}
	return nil
}

//handle SELECT DATABASE()
func (mce *MysqlCmdExecutor) handleSelectDatabase(sel *tree.Select) error{
	var err error = nil
	ses := mce.routine.GetSession()
	proto := mce.routine.GetClientProtocol().(*MysqlProtocol)

	col := new(MysqlColumn)
	col.SetName("DATABASE()")
	col.SetColumnType(defines.MYSQL_TYPE_VARCHAR)
	ses.Mrs.AddColumn(col)
	val := mce.routine.db
	if val == "" {
		val = "NULL"
	}
	ses.Mrs.AddRow([]interface{}{val})

	mer := NewMysqlExecutionResult(0, 0, 0, 0, ses.Mrs)
	resp := NewResponse(ResultResponse, 0, int(COM_QUERY), mer)

	if err = proto.SendResponse(resp); err != nil {
		return fmt.Errorf("routine send response failed. error:%v ", err)
	}
	return nil
}

//execute query
func (mce *MysqlCmdExecutor) doComQuery(sql string) error {
	ses := mce.routine.GetSession()
	proto := mce.routine.GetClientProtocol().(*MysqlProtocol)
	pdHook := mce.routine.GetPDCallback().(*PDCallbackImpl)
	statementCount := uint64(1)

	//pin the epoch with 1
	epoch, _ := pdHook.IncQueryCountAtCurrentEpoch(statementCount)
	defer func() {
		ep, stmtCnt := pdHook.DecQueryCountAtEpoch(epoch, statementCount)
		if ep != epoch || stmtCnt != 0 {
			panic(fmt.Errorf("statement_count needs zero, but actually it is %d at epoch %d \n", stmtCnt, ep))
		}
	}()

	proc := process.New(ses.GuestMmu, ses.Mempool)
	proc.Id = mce.getNextProcessId()
	proc.Lim.Size = ses.Pu.SV.GetProcessLimitationSize()
	proc.Lim.BatchRows = ses.Pu.SV.GetProcessLimitationBatchRows()
	proc.Lim.PartitionRows = ses.Pu.SV.GetProcessLimitationPartitionRows()
	proc.Refer = make(map[string]uint64)

	comp := compile.New(mce.routine.db, sql, mce.routine.user, ses.Pu.StorageEngine, ses.Pu.ClusterNodes, proc)
	execs, err := comp.Compile()
	if err != nil {
		return err
	}

	var choose = !ses.Pu.SV.GetSendRow()

	ses.Mrs = &MysqlResultSet{}

	defer func() {
		ses.Mrs = nil
	}()

	for _, exec := range execs {
		stmt := exec.Statement()

		//temp try 0 epoch
		pdHook.IncQueryCountAtEpoch(epoch, 1)
		statementCount++

		switch st := stmt.(type) {
		case *tree.Select:
			if sc,ok := st.Select.(*tree.SelectClause) ; ok {
				if len(sc.Exprs) == 1 {
					if fe,ok := sc.Exprs[0].Expr.(*tree.FuncExpr); ok {
						if un,ok := fe.Func.FunctionReference.(*tree.UnresolvedName); ok {
							if strings.ToUpper(un.Parts[0]) == "DATABASE" {
								err = mce.handleSelectDatabase(st)
								if err != nil{
									return err
								}

								//next statement
								continue
							}
						}
					}
				}
			}
		}

		//check database
		if mce.routine.db == "" {
			//if none database has been selected, database operations must be failed.
			switch stmt.(type) {
			case *tree.ShowDatabases, *tree.CreateDatabase, *tree.ShowWarnings, *tree.ShowErrors,
				*tree.ShowStatus,*tree.DropDatabase:
			default:
				return NewMysqlError(ER_NO_DB_ERROR)
			}
		}

		var selfHandle = false

		switch st := stmt.(type) {
		case *tree.Use:
			selfHandle = true
			resp, err := mce.routine.ChangeDB(st.Name)
			if err != nil {
				return err
			}
			if resp != nil {
				err := proto.SendResponse(resp)
				if err != nil {
					return err
				}
			}
		}

		if selfHandle {
			continue
		}
		if err = exec.SetSchema(mce.routine.db); err != nil {
			return err
		}

		cmpBegin := time.Now()
		if err = exec.Compile(mce.routine, getDataFromPipeline); err != nil {
			return err
		}

		if ses.Pu.SV.GetRecordTimeElapsedOfSqlRequest() {
			fmt.Printf("time of Exec.Compile : %s \n",time.Since(cmpBegin).String())
		}

		switch stmt.(type) {
		//produce result set
		case *tree.Select,
			*tree.ShowCreate, *tree.ShowCreateDatabase, *tree.ShowTables, *tree.ShowDatabases, *tree.ShowColumns,
			*tree.ShowProcessList, *tree.ShowErrors, *tree.ShowWarnings, *tree.ShowVariables, *tree.ShowStatus,
			*tree.ShowIndex,
			*tree.ExplainFor, *tree.ExplainAnalyze, *tree.ExplainStmt:
			columns := exec.Columns()
			if choose {

				/*
					Step 1 : send column count and column definition.
				*/
				//send column count
				colCnt := uint64(len(columns))
				err := proto.SendColumnCountPacket(colCnt)
				if err != nil {
					return err
				}
				//send columns
				//column_count * Protocol::ColumnDefinition packets
				cmd := ses.Cmd
				for _, c := range columns {
					col := new(MysqlColumn)
					col.SetName(c.Name)
					switch c.Typ {
					case types.T_int8:
						col.SetColumnType(defines.MYSQL_TYPE_TINY)
					case types.T_uint8:
						col.SetColumnType(defines.MYSQL_TYPE_TINY)
						col.SetSigned(true)
					case types.T_int16:
						col.SetColumnType(defines.MYSQL_TYPE_SHORT)
					case types.T_uint16:
						col.SetColumnType(defines.MYSQL_TYPE_SHORT)
						col.SetSigned(true)
					case types.T_int32:
						col.SetColumnType(defines.MYSQL_TYPE_LONG)
					case types.T_uint32:
						col.SetColumnType(defines.MYSQL_TYPE_LONG)
						col.SetSigned(true)
					case types.T_int64:
						col.SetColumnType(defines.MYSQL_TYPE_LONGLONG)
					case types.T_uint64:
						col.SetColumnType(defines.MYSQL_TYPE_LONGLONG)
						col.SetSigned(true)
					case types.T_float32:
						col.SetColumnType(defines.MYSQL_TYPE_FLOAT)
					case types.T_float64:
						col.SetColumnType(defines.MYSQL_TYPE_DOUBLE)
					case types.T_char:
						col.SetColumnType(defines.MYSQL_TYPE_STRING)
					case types.T_varchar:
						col.SetColumnType(defines.MYSQL_TYPE_VARCHAR)
					default:
						return fmt.Errorf("RunWhileSend : unsupported type %d \n", c.Typ)
					}
					ses.Mrs.AddColumn(col)

					//fmt.Printf("doComQuery col name %v type %v \n",col.Name(),col.ColumnType())
					/*
						mysql COM_QUERY response: send the column definition per column
					*/
					err := proto.SendColumnDefinitionPacket(col, cmd)
					if err != nil {
						return err
					}
				}

				/*
					mysql COM_QUERY response: End after the column has been sent.
					send EOF packet
				*/
				err = proto.SendEOFPacketIf(0, 0)
				if err != nil {
					return err
				}

				runBegin := time.Now()
				/*
					Step 2: Start pipeline
					Producing the data row and sending the data row
				*/
				if er := exec.Run(epoch); er != nil {
					return er
				}
				if ses.Pu.SV.GetRecordTimeElapsedOfSqlRequest() {
					fmt.Printf("time of Exec.Run : %s \n",time.Since(runBegin).String())
				}
				/*
					Step 3: Say goodbye
					mysql COM_QUERY response: End after the data row has been sent.
					After all row data has been sent, it sends the EOF or OK packet.
				*/
				err = proto.sendEOFOrOkPacket(0, 0)
				if err != nil {
					return err
				}
			} else {
				for _, c := range columns {
					col := new(MysqlColumn)
					col.SetName(c.Name)
					switch c.Typ {
					case types.T_int8:
						col.SetColumnType(defines.MYSQL_TYPE_TINY)
					case types.T_uint8:
						col.SetColumnType(defines.MYSQL_TYPE_TINY)
						col.SetSigned(true)
					case types.T_int16:
						col.SetColumnType(defines.MYSQL_TYPE_SHORT)
					case types.T_uint16:
						col.SetColumnType(defines.MYSQL_TYPE_SHORT)
						col.SetSigned(true)
					case types.T_int32:
						col.SetColumnType(defines.MYSQL_TYPE_LONG)
					case types.T_uint32:
						col.SetColumnType(defines.MYSQL_TYPE_LONG)
						col.SetSigned(true)
					case types.T_int64:
						col.SetColumnType(defines.MYSQL_TYPE_LONGLONG)
					case types.T_uint64:
						col.SetColumnType(defines.MYSQL_TYPE_LONGLONG)
						col.SetSigned(true)
					case types.T_float32:
						col.SetColumnType(defines.MYSQL_TYPE_FLOAT)
					case types.T_float64:
						col.SetColumnType(defines.MYSQL_TYPE_DOUBLE)
					case types.T_char:
						col.SetColumnType(defines.MYSQL_TYPE_STRING)
					case types.T_varchar:
						col.SetColumnType(defines.MYSQL_TYPE_VARCHAR)
					default:
						return fmt.Errorf("RunWhileSend : unsupported type %d \n", c.Typ)
					}

					ses.Mrs.AddColumn(col)
				}

				if er := exec.Run(epoch); er != nil {
					return er
				}

				mer := NewMysqlExecutionResult(0, 0, 0, 0, ses.Mrs)
				resp := NewResponse(ResultResponse, 0, int(COM_QUERY), mer)

				if err = proto.SendResponse(resp); err != nil {
					return fmt.Errorf("routine send response failed. error:%v ", err)
				}
			}
		//just status, no result set
		case *tree.CreateTable, *tree.DropTable, *tree.CreateDatabase, *tree.DropDatabase,
			*tree.CreateIndex, *tree.DropIndex,
			*tree.Insert, *tree.Delete, *tree.Update,
			*tree.BeginTransaction, *tree.CommitTransaction, *tree.RollbackTransaction,
			*tree.SetVar,
			*tree.Load,
			*tree.CreateUser, *tree.DropUser, *tree.AlterUser,
			*tree.CreateRole, *tree.DropRole,
			*tree.Revoke, *tree.Grant,
			*tree.SetDefaultRole, *tree.SetRole, *tree.SetPassword:
			/*
				Step 1: Start
			*/
			if er := exec.Run(epoch); er != nil {
				return er
			}

			//record ddl drop xxx after the success
			switch stmt.(type) {
			case *tree.DropTable, *tree.DropDatabase,
					*tree.DropIndex, *tree.DropUser, *tree.DropRole:
				//test ddl
				pdHook.IncDDLCountAtEpoch(epoch,1)
			}

			/*
				Step 2: Echo client
			*/
			resp := NewResponse(
				OkResponse,
				0,
				int(COM_QUERY),
				nil,
			)
			if err = proto.SendResponse(resp); err != nil {
				return err
			}
		}
	}

	return nil
}

// ExecRequest the server execute the commands from the client following the mysql's routine
func (mce *MysqlCmdExecutor) ExecRequest(req *Request) (*Response, error) {
	var resp *Response = nil
	fmt.Printf("cmd %v \n", req.GetCmd())

	ses := mce.routine.GetSession()
	if ses.Pu.SV.GetRejectWhenHeartbeatFromPDLeaderIsTimeout() {
		pdHook := mce.routine.GetPDCallback().(*PDCallbackImpl)
		if !pdHook.CanAcceptSomething() {
			resp = NewResponse(
				ErrorResponse,
				0,
				req.GetCmd(),
				fmt.Errorf("heartbeat from pdleader is timeout. the server reject sql request. cmd %d \n", req.GetCmd()),
			)
			return resp, nil
		}
	}

	switch uint8(req.GetCmd()) {
	case COM_QUIT:
		resp = NewResponse(
			OkResponse,
			0,
			int(COM_QUIT),
			nil,
		)
		return resp, nil
	case COM_QUERY:
		var query = string(req.GetData().([]byte))
		mce.addSqlCount(1)
		fmt.Printf("query:%s \n", query)
		err := mce.doComQuery(query)
		if err != nil {
			resp = NewResponse(
				ErrorResponse,
				0,
				int(COM_QUERY),
				err,
			)
		}
		return resp, nil
	case COM_INIT_DB:
		var dbname = string(req.GetData().([]byte))
		resp, err := mce.routine.ChangeDB(dbname)
		if err != nil {
			resp = NewResponse(
				ErrorResponse,
				0,
				int(COM_INIT_DB),
				err,
			)
		} else {
			err := mce.routine.protocol.SendResponse(resp)
			if err != nil {
				return nil, err
			}
		}

		return nil, nil
	default:
		err := fmt.Errorf("unsupported command. 0x%x \n", req.GetCmd())
		resp = NewResponse(
			ErrorResponse,
			0,
			req.GetCmd(),
			err,
		)
	}
	return resp, nil
}

func (mce *MysqlCmdExecutor) Close() {
	//TODO:
}

func NewMysqlCmdExecutor() *MysqlCmdExecutor {
	return &MysqlCmdExecutor{}
}