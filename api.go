package pick

type SelectorMustAPI struct {
	*Picker
	onErr []func(selector string, err error)
}

func (a SelectorMustAPI) Bool(selector string) bool {
	return SelectorMust(a.Picker, selector, a.caster.AsBool, a.onErr...)
}

func (a SelectorMustAPI) BoolSlice(selector string) []bool {
	return SelectorMust(a.Picker, selector, a.caster.AsBoolSlice, a.onErr...)
}

func (a SelectorMustAPI) Byte(selector string) byte {
	return SelectorMust(a.Picker, selector, a.caster.AsByte, a.onErr...)
}

func (a SelectorMustAPI) ByteSlice(selector string) []byte {
	return SelectorMust(a.Picker, selector, a.caster.AsByteSlice, a.onErr...)
}

func (a SelectorMustAPI) Float32(selector string) float32 {
	return SelectorMust(a.Picker, selector, a.caster.AsFloat32, a.onErr...)
}

func (a SelectorMustAPI) Float32Slice(selector string) []float32 {
	return SelectorMust(a.Picker, selector, a.caster.AsFloat32Slice, a.onErr...)
}

func (a SelectorMustAPI) Float64(selector string) float64 {
	return SelectorMust(a.Picker, selector, a.caster.AsFloat64, a.onErr...)
}

func (a SelectorMustAPI) Float64Slice(selector string) []float64 {
	return SelectorMust(a.Picker, selector, a.caster.AsFloat64Slice, a.onErr...)
}

func (a SelectorMustAPI) Int(selector string) int {
	return SelectorMust(a.Picker, selector, a.caster.AsInt, a.onErr...)
}

func (a SelectorMustAPI) IntSlice(selector string) []int {
	return SelectorMust(a.Picker, selector, a.caster.AsIntSlice, a.onErr...)
}

func (a SelectorMustAPI) Int8(selector string) int8 {
	return SelectorMust(a.Picker, selector, a.caster.AsInt8, a.onErr...)
}

func (a SelectorMustAPI) Int8Slice(selector string) []int8 {
	return SelectorMust(a.Picker, selector, a.caster.AsInt8Slice, a.onErr...)
}

func (a SelectorMustAPI) Int16(selector string) int16 {
	return SelectorMust(a.Picker, selector, a.caster.AsInt16, a.onErr...)
}

func (a SelectorMustAPI) Int16Slice(selector string) []int16 {
	return SelectorMust(a.Picker, selector, a.caster.AsInt16Slice, a.onErr...)
}

func (a SelectorMustAPI) Int32(selector string) int32 {
	return SelectorMust(a.Picker, selector, a.caster.AsInt32, a.onErr...)
}

func (a SelectorMustAPI) Int32Slice(selector string) []int32 {
	return SelectorMust(a.Picker, selector, a.caster.AsInt32Slice, a.onErr...)
}

func (a SelectorMustAPI) Int64(selector string) int64 {
	return SelectorMust(a.Picker, selector, a.caster.AsInt64, a.onErr...)
}

func (a SelectorMustAPI) Int64Slice(selector string) []int64 {
	return SelectorMust(a.Picker, selector, a.caster.AsInt64Slice, a.onErr...)
}

func (a SelectorMustAPI) Uint(selector string) uint {
	return SelectorMust(a.Picker, selector, a.caster.AsUint, a.onErr...)
}

func (a SelectorMustAPI) UintSlice(selector string) []uint {
	return SelectorMust(a.Picker, selector, a.caster.AsUintSlice, a.onErr...)
}

func (a SelectorMustAPI) Uint8(selector string) uint8 {
	return SelectorMust(a.Picker, selector, a.caster.AsUint8, a.onErr...)
}

func (a SelectorMustAPI) Uint8Slice(selector string) []uint8 {
	return SelectorMust(a.Picker, selector, a.caster.AsUint8Slice, a.onErr...)
}

func (a SelectorMustAPI) Uint16(selector string) uint16 {
	return SelectorMust(a.Picker, selector, a.caster.AsUint16, a.onErr...)
}

func (a SelectorMustAPI) Uint16Slice(selector string) []uint16 {
	return SelectorMust(a.Picker, selector, a.caster.AsUint16Slice, a.onErr...)
}

func (a SelectorMustAPI) Uint32(selector string) uint32 {
	return SelectorMust(a.Picker, selector, a.caster.AsUint32, a.onErr...)
}

func (a SelectorMustAPI) Uint32Slice(selector string) []uint32 {
	return SelectorMust(a.Picker, selector, a.caster.AsUint32Slice, a.onErr...)
}

func (a SelectorMustAPI) Uint64(selector string) uint64 {
	return SelectorMust(a.Picker, selector, a.caster.AsUint64, a.onErr...)
}

func (a SelectorMustAPI) Uint64Slice(selector string) []uint64 {
	return SelectorMust(a.Picker, selector, a.caster.AsUint64Slice, a.onErr...)
}

func (a SelectorMustAPI) String(selector string) string {
	return SelectorMust(a.Picker, selector, a.caster.AsString, a.onErr...)
}

func (a SelectorMustAPI) StringSlice(selector string) []string {
	return SelectorMust(a.Picker, selector, a.caster.AsStringSlice, a.onErr...)
}

type PathAPI struct {
	*Picker
}

func (a PathAPI) Bool(path ...Key) (bool, error) {
	return Path(a.Picker, path, a.caster.AsBool)
}

func (a PathAPI) BoolSlice(path ...Key) ([]bool, error) {
	return Path(a.Picker, path, a.caster.AsBoolSlice)
}

func (a PathAPI) Byte(path ...Key) (byte, error) {
	return Path(a.Picker, path, a.caster.AsByte)
}

func (a PathAPI) ByteSlice(path ...Key) ([]byte, error) {
	return Path(a.Picker, path, a.caster.AsByteSlice)
}

func (a PathAPI) Float32(path ...Key) (float32, error) {
	return Path(a.Picker, path, a.caster.AsFloat32)
}

func (a PathAPI) Float32Slice(path ...Key) ([]float32, error) {
	return Path(a.Picker, path, a.caster.AsFloat32Slice)
}

func (a PathAPI) Float64(path ...Key) (float64, error) {
	return Path(a.Picker, path, a.caster.AsFloat64)
}

func (a PathAPI) Float64Slice(path ...Key) ([]float64, error) {
	return Path(a.Picker, path, a.caster.AsFloat64Slice)
}

func (a PathAPI) Int(path ...Key) (int, error) {
	return Path(a.Picker, path, a.caster.AsInt)
}

func (a PathAPI) IntSlice(path ...Key) ([]int, error) {
	return Path(a.Picker, path, a.caster.AsIntSlice)
}

func (a PathAPI) Int8(path ...Key) (int8, error) {
	return Path(a.Picker, path, a.caster.AsInt8)
}

func (a PathAPI) Int8Slice(path ...Key) ([]int8, error) {
	return Path(a.Picker, path, a.caster.AsInt8Slice)
}

func (a PathAPI) Int16(path ...Key) (int16, error) {
	return Path(a.Picker, path, a.caster.AsInt16)
}

func (a PathAPI) Int16Slice(path ...Key) ([]int16, error) {
	return Path(a.Picker, path, a.caster.AsInt16Slice)
}

func (a PathAPI) Int32(path ...Key) (int32, error) {
	return Path(a.Picker, path, a.caster.AsInt32)
}

func (a PathAPI) Int32Slice(path ...Key) ([]int32, error) {
	return Path(a.Picker, path, a.caster.AsInt32Slice)
}

func (a PathAPI) Int64(path ...Key) (int64, error) {
	return Path(a.Picker, path, a.caster.AsInt64)
}

func (a PathAPI) Int64Slice(path ...Key) ([]int64, error) {
	return Path(a.Picker, path, a.caster.AsInt64Slice)
}

func (a PathAPI) Uint(path ...Key) (uint, error) {
	return Path(a.Picker, path, a.caster.AsUint)
}

func (a PathAPI) UintSlice(path ...Key) ([]uint, error) {
	return Path(a.Picker, path, a.caster.AsUintSlice)
}

func (a PathAPI) Uint8(path ...Key) (uint8, error) {
	return Path(a.Picker, path, a.caster.AsUint8)
}

func (a PathAPI) Uint8Slice(path ...Key) ([]uint8, error) {
	return Path(a.Picker, path, a.caster.AsUint8Slice)
}

func (a PathAPI) Uint16(path ...Key) (uint16, error) {
	return Path(a.Picker, path, a.caster.AsUint16)
}

func (a PathAPI) Uint16Slice(path ...Key) ([]uint16, error) {
	return Path(a.Picker, path, a.caster.AsUint16Slice)
}

func (a PathAPI) Uint32(path ...Key) (uint32, error) {
	return Path(a.Picker, path, a.caster.AsUint32)
}

func (a PathAPI) Uint32Slice(path ...Key) ([]uint32, error) {
	return Path(a.Picker, path, a.caster.AsUint32Slice)
}

func (a PathAPI) Uint64(path ...Key) (uint64, error) {
	return Path(a.Picker, path, a.caster.AsUint64)
}

func (a PathAPI) Uint64Slice(path ...Key) ([]uint64, error) {
	return Path(a.Picker, path, a.caster.AsUint64Slice)
}

func (a PathAPI) String(path ...Key) (string, error) {
	return Path(a.Picker, path, a.caster.AsString)
}

func (a PathAPI) StringSlice(path ...Key) ([]string, error) {
	return Path(a.Picker, path, a.caster.AsStringSlice)
}

type PathMustAPI struct {
	*Picker
	onErr []func(selector string, err error)
}

func (a PathMustAPI) Bool(path ...Key) bool {
	return PathMust(a.Picker, path, a.caster.AsBool, a.onErr...)
}

func (a PathMustAPI) BoolSlice(path ...Key) []bool {
	return PathMust(a.Picker, path, a.caster.AsBoolSlice, a.onErr...)
}

func (a PathMustAPI) Byte(path ...Key) byte {
	return PathMust(a.Picker, path, a.caster.AsByte, a.onErr...)
}

func (a PathMustAPI) ByteSlice(path ...Key) []byte {
	return PathMust(a.Picker, path, a.caster.AsByteSlice, a.onErr...)
}

func (a PathMustAPI) Float32(path ...Key) float32 {
	return PathMust(a.Picker, path, a.caster.AsFloat32, a.onErr...)
}

func (a PathMustAPI) Float32Slice(path ...Key) []float32 {
	return PathMust(a.Picker, path, a.caster.AsFloat32Slice, a.onErr...)
}

func (a PathMustAPI) Float64(path ...Key) float64 {
	return PathMust(a.Picker, path, a.caster.AsFloat64, a.onErr...)
}

func (a PathMustAPI) Float64Slice(path ...Key) []float64 {
	return PathMust(a.Picker, path, a.caster.AsFloat64Slice, a.onErr...)
}

func (a PathMustAPI) Int(path ...Key) int {
	return PathMust(a.Picker, path, a.caster.AsInt, a.onErr...)
}

func (a PathMustAPI) IntSlice(path ...Key) []int {
	return PathMust(a.Picker, path, a.caster.AsIntSlice, a.onErr...)
}

func (a PathMustAPI) Int8(path ...Key) int8 {
	return PathMust(a.Picker, path, a.caster.AsInt8, a.onErr...)
}

func (a PathMustAPI) Int8Slice(path ...Key) []int8 {
	return PathMust(a.Picker, path, a.caster.AsInt8Slice, a.onErr...)
}

func (a PathMustAPI) Int16(path ...Key) int16 {
	return PathMust(a.Picker, path, a.caster.AsInt16, a.onErr...)
}

func (a PathMustAPI) Int16Slice(path ...Key) []int16 {
	return PathMust(a.Picker, path, a.caster.AsInt16Slice, a.onErr...)
}

func (a PathMustAPI) Int32(path ...Key) int32 {
	return PathMust(a.Picker, path, a.caster.AsInt32, a.onErr...)
}

func (a PathMustAPI) Int32Slice(path ...Key) []int32 {
	return PathMust(a.Picker, path, a.caster.AsInt32Slice, a.onErr...)
}

func (a PathMustAPI) Int64(path ...Key) int64 {
	return PathMust(a.Picker, path, a.caster.AsInt64, a.onErr...)
}

func (a PathMustAPI) Int64Slice(path ...Key) []int64 {
	return PathMust(a.Picker, path, a.caster.AsInt64Slice, a.onErr...)
}

func (a PathMustAPI) Uint(path ...Key) uint {
	return PathMust(a.Picker, path, a.caster.AsUint, a.onErr...)
}

func (a PathMustAPI) UintSlice(path ...Key) []uint {
	return PathMust(a.Picker, path, a.caster.AsUintSlice, a.onErr...)
}

func (a PathMustAPI) Uint8(path ...Key) uint8 {
	return PathMust(a.Picker, path, a.caster.AsUint8, a.onErr...)
}

func (a PathMustAPI) Uint8Slice(path ...Key) []uint8 {
	return PathMust(a.Picker, path, a.caster.AsUint8Slice, a.onErr...)
}

func (a PathMustAPI) Uint16(path ...Key) uint16 {
	return PathMust(a.Picker, path, a.caster.AsUint16, a.onErr...)
}

func (a PathMustAPI) Uint16Slice(path ...Key) []uint16 {
	return PathMust(a.Picker, path, a.caster.AsUint16Slice, a.onErr...)
}

func (a PathMustAPI) Uint32(path ...Key) uint32 {
	return PathMust(a.Picker, path, a.caster.AsUint32, a.onErr...)
}

func (a PathMustAPI) Uint32Slice(path ...Key) []uint32 {
	return PathMust(a.Picker, path, a.caster.AsUint32Slice, a.onErr...)
}

func (a PathMustAPI) Uint64(path ...Key) uint64 {
	return PathMust(a.Picker, path, a.caster.AsUint64, a.onErr...)
}

func (a PathMustAPI) Uint64Slice(path ...Key) []uint64 {
	return PathMust(a.Picker, path, a.caster.AsUint64Slice, a.onErr...)
}

func (a PathMustAPI) String(path ...Key) string {
	return PathMust(a.Picker, path, a.caster.AsString, a.onErr...)
}

func (a PathMustAPI) StringSlice(path ...Key) []string {
	return PathMust(a.Picker, path, a.caster.AsStringSlice, a.onErr...)
}
