package main

/*
compression is a big topic but it does boil down to the
trade-off between access speed and compression ratio

It also matters what kind of database this is, analytical databases tend
to store alot more data and compression is a little more amenable to the workload.
*/

// there's not much to show/code here as most of the work goes on in the underlying library ie zstd or snappy.

// --snipped from pebble
// https://github.com/cockroachdb/pebble/blob/4f4be4295c9c0adef551093c4b349112f686cc61/sstable/compression.go#L84
/*
// compressBlock compresses an SST block, using compressBuf as the desired destination.
func compressBlock(
	compression Compression, b []byte, compressedBuf []byte,
) (blockType blockType, compressed []byte) {
	switch compression {
	case SnappyCompression:
		return snappyCompressionBlockType, snappy.Encode(compressedBuf, b)
	case NoCompression:
		return noCompressionBlockType, b
	}

	if len(compressedBuf) < binary.MaxVarintLen64 {
		compressedBuf = append(compressedBuf, make([]byte, binary.MaxVarintLen64-len(compressedBuf))...)
	}
	varIntLen := binary.PutUvarint(compressedBuf, uint64(len(b)))
	switch compression {
	case ZstdCompression:
		return zstdCompressionBlockType, encodeZstd(compressedBuf, varIntLen, b)
	default:
		retu
*/
