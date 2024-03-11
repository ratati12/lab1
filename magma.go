package main

func generate_keys(key [8]uint32) [32]uint32 {
	var iter_keys [32]uint32
	for i := 0; i < 32; i++ {
		if i < 24 {
			iter_keys[i] = key[i%8]
		}
		if i >= 24 {
			iter_keys[i] = key[7-i%8]
		}
	}
	key = [8]uint32{0x0}
	return iter_keys
}

func t(input uint32) uint32 {
	pi := [8][16]byte{
		{12, 4, 6, 2, 10, 5, 11, 9, 14, 8, 13, 7, 0, 3, 15, 1},
		{6, 8, 2, 3, 9, 10, 5, 12, 1, 14, 4, 7, 11, 13, 0, 15},
		{11, 3, 5, 8, 2, 15, 10, 13, 14, 1, 7, 4, 12, 9, 6, 0},
		{12, 8, 2, 1, 13, 4, 15, 6, 7, 0, 10, 5, 3, 14, 9, 11},
		{7, 15, 5, 10, 8, 1, 6, 13, 0, 9, 3, 14, 11, 4, 2, 12},
		{5, 13, 15, 6, 9, 2, 12, 10, 11, 7, 8, 1, 4, 3, 14, 0},
		{8, 14, 2, 5, 6, 9, 1, 12, 15, 4, 11, 0, 13, 10, 3, 7},
		{1, 7, 14, 13, 0, 5, 8, 3, 4, 15, 10, 6, 9, 12, 11, 2}}
	var a [4]uint8
	var output uint32

	for i := 3; i >= 0; i-- {
		a[i] = uint8(input >> (8 * i))
		fp := uint8(pi[i*2][a[i]&0x0f])
		sc := uint8(pi[i*2+1][(a[i]&0xf0)>>4])
		output = output << 8
		output = output + uint32((sc<<4)|fp)
	}
	return output
}

func uint64ToUint32Array(value uint64) [2]uint32 {
	// Преобразование uint64 в два uint32
	lower32 := uint32(value)
	upper32 := uint32(value >> 32)

	// Создание и возвращение массива [2]uint32
	return [2]uint32{upper32, lower32}
}

func uint32ArrayToUint64(array [2]uint32) uint64 {
	// Соединение двух uint32 в один uint64
	result := uint64(array[1]) | (uint64(array[0]) << 32)
	return result
}

func rol32(n uint32, nPos int) uint32 {
	return ((n << nPos) & (1<<32 - 1)) | ((n >> (32 - nPos)) & (1<<32 - 1))
}

func g(a uint32, k uint32) uint32 {
	return rol32(t(a+k), 11)
}
func big_g(a uint32, b uint32) uint32 {
	return a ^ b
}

func encrypt(block64 uint64, key [8]uint32) uint64 {
	block32 := [2]uint32(uint64ToUint32Array(block64))
	k := [32]uint32(generate_keys(key))
	key = [8]uint32{0x0}
	for i := 0; i < 32; i++ {
		temp := uint32(block32[1])
		block32[1] = big_g(block32[0], g(block32[1], k[i]))
		block32[0] = temp
	}
	temp := block32[1]
	block32[1] = block32[0]
	block32[0] = temp
	block64 = uint32ArrayToUint64(block32)
	return block64
}

func decrypt(block64 uint64, key [8]uint32) uint64 {
	block32 := [2]uint32(uint64ToUint32Array(block64))
	k := [32]uint32(generate_keys(key))
	key = [8]uint32{0x0}
	for i := 31; i >= 0; i-- {
		temp := uint32(block32[1])
		block32[1] = big_g(block32[0], g(block32[1], k[i]))
		block32[0] = temp
	}
	temp := block32[1]
	block32[1] = block32[0]
	block32[0] = temp
	block64 = uint32ArrayToUint64(block32)
	return block64
}
