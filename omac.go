package main

func rol64(n uint64, nPos int) uint64 {
	return ((n << nPos) & (1<<64 - 1)) | ((n >> (64 - nPos)) & (1<<64 - 1))
}

func msb(value uint64) uint64 {
	return (value >> 63) & 0x01
}

func omac_subkey(key [8]uint32) [2]uint64 {
	k0 := uint64(encrypt(0x0, key))
	key = [8]uint32{0x0}
	var k [2]uint64
	if msb(k0) == 0 {
		k[0] = rol64(k0, 1)
	} else {
		k[0] = rol64(k0, 1) ^ 0x1b
	}

	if msb(k[0]) == 0 {
		k[1] = rol64(k[0], 1)
	} else {
		k[1] = rol64(k[0], 1) ^ 0x1b
	}
	return k
}

func omac(block []uint64, key [8]uint32, n int64) uint64 {
	omacValue := uint64(0x0)
	omac_subkeys := [2]uint64(omac_subkey(key))

	for i := int64(0); i < n-1; i++ {
		// XOR с предыдущим значением
		omacValue ^= block[i]
		// Шифрование значения
		omacValue = encrypt(omacValue, key)
	}
	omacValue ^= block[n-1]
	// Дополнительный шаг для последнего блока
	omacValue ^= omac_subkeys[0]

	// Шифрование окончательного значения
	return encrypt(omacValue, key) >> 32
}
