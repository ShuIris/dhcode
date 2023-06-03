package main

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
)

type Node struct {
	nodeID  string
	buckets []*Bucket
	keys    map[string][]byte
}

type Bucket struct {
	ids []string
}

var nodesMap map[string]*Node

var SetNodes []string
var GetNodes []string

func (s *Node) FindNode(nodeID string, array []string) []string {
	var nodes []string
	var return_node []string
	if s.nodeID == nodeID {
		nodes = append(nodes, s.nodeID, s.nodeID)
		return nodes
	}

	if len(s.buckets) == 0 {
		s.InsertNode(nodeID)
		return array
	}

	//寻找到对应的桶
	result := findBucket(s.nodeID, nodeID)
	var bucket *Bucket
	if result >= (len(s.buckets) - 1) {
		bucket = s.buckets[len(s.buckets)-1]
	} else {
		bucket = s.buckets[result]
	}

	//判断桶中是否存在该节点
	for _, v := range bucket.ids {
		if v == nodeID {
			nodes = append(nodes, v, v)
			return nodes
		}
	}

	var node1, node2 string
	var nodeNum int
	//不存在就进行递归,桶中选取随机的两个节点
	if len(bucket.ids) == 2 {
		node1 = bucket.ids[0]
		node2 = bucket.ids[1]
		nodeNum = 2
	} else if len(bucket.ids) == 1 {
		node1 = bucket.ids[0]
		nodeNum = 1
	} else if len(bucket.ids) > 2 {
		index1, index2 := GetRandom2()
		node1 = bucket.ids[index1]
		node2 = bucket.ids[index2]
		nodeNum = 2
	}


	if nodeNum == 2 {
		//在第一遍比较的时候，对于array中已经发生交换的元素，在第二次比较的时候就会进行跳过
		isUpdate := -1
		for i, v := range array {
			result := compareGetMin(nodeID, v, node1)
			if result == node1 {
				array[i] = node1
				isUpdate = i
				return_node = append(return_node, nodesMap[node1].FindNode(nodeID, array)...)
			}
		}

		for i := len(array) - 1; i >= 0; i-- {
			if i == isUpdate {
				continue
			}
			result := compareGetMin(nodeID, array[i], node2)
			if result == node2 {
				array[i] = node2
				return_node = append(return_node, nodesMap[node2].FindNode(nodeID, array)...)
			}
		}
	} else if nodeNum == 1 {
		num := new(big.Int)
		num1 := new(big.Int)
		num2 := new(big.Int)
		num3 := new(big.Int)
		num.SetString(nodeID, 2)
		num1.SetString(array[0], 2)
		num2.SetString(array[1], 2)
		num3.SetString(node1, 2)
		//选出array中最大的与新结点进行比较
		result1 := new(big.Int)
		result1.Xor(num, num1)
		result2 := new(big.Int)
		result2.Xor(num, num2)

		if result1.Cmp(result2) > 0 {
			if num3.Cmp(num1) < 0 {
				array[0] = node1
			}
			return_node = append(return_node, nodesMap[node1].FindNode(nodeID, array)...)
		} else if result1.Cmp(result2) < 0 {
			if num3.Cmp(num2) < 0 {
				array[1] = node1
			}
			return_node = append(return_node, nodesMap[node1].FindNode(nodeID, array)...)
		}
	} else {
		// s.InsertNode(nodeID)
		return array
	}

	//将新节点FindNode返回的列表中的节点与传入的列表中的节点进行比对，选出最近的两个节点进行返回
	//将返回的节点信息中与array中不相同的节点加入array中
	for _, v := range return_node {
		if v == array[0] || v == array[1] {
			continue
		}
		array = append(array, v)
	}

	//如果array的长度小于等于2的话，说明找到最小的2个节点或1一个节点的信息
	if len(array) <= 2 {
		// s.InsertNode(nodeID)
		return array
	}
	//否则就挑选出arry中距离最小的两个元素放入nodes中
	index_min1, index_min2 := get2MinIndex(array, nodeID)
	nodes = append(nodes, array[index_min1], array[index_min2])

	// s.InsertNode(nodeID)
	return nodes
}

func compareGetMin(targetValue, value1, value2 string) string {
	num := new(big.Int)
	num1 := new(big.Int)
	num2 := new(big.Int)
	num.SetString(targetValue, 2)
	num1.SetString(value1, 2)
	num2.SetString(value2, 2)
	//计算出距离
	result1 := new(big.Int)
	result1.Xor(num, num1)
	result2 := new(big.Int)
	result2.Xor(num, num2)

	if result1.Cmp(result2) < 0 {
		return value1
	} else {
		return value2
	}
}

func get2MinIndex(nodes []string, targetNode string) (int, int) {
	var distance []*big.Int
	for _, v := range nodes {
		num1 := new(big.Int)
		num2 := new(big.Int)
		num1.SetString(targetNode, 2)
		num2.SetString(v, 2)
		result := new(big.Int)
		result.Xor(num1, num2)
		distance = append(distance, result)
	}

	//找出aray中最近的两个元素
	temp1 := distance[0]
	index_min1 := 0
	temp2 := distance[1]
	index_min2 := 1
	//先找到第一个最近的
	for i := 1; i < len(distance); i++ {
		if distance[i].Cmp(temp1) < 0 {
			temp1 = distance[i]
			index_min1 = i
		}
	}
	//在找到第二个近的
	for i, v := range distance {
		if i == index_min1 || i == 1 {
			continue
		}
		if v.Cmp(temp2) < 0 {
			temp2 = distance[i]
			index_min2 = i
		}
	}
	return index_min1, index_min2
}

// 生成两个随机数，0~2之间
func GetRandom2() (int, int) {
	var nums [2]int
	// 随机生成两个不重复的整数
	for i := range nums {
		num, err := rand.Int(rand.Reader, big.NewInt(3))
		if err != nil {
			// 处理错误
			return -1, -1
		}
		nums[i] = int(num.Int64())
	}
	for nums[0] == nums[1] {
		num, err := rand.Int(rand.Reader, big.NewInt(3))
		if err != nil {
			// 处理错误
			fmt.Println(err)
			return -1, -1
		}
		nums[1] = int(num.Int64())
	}
	return nums[0], nums[1]
}



func (s *Node) InsertNode(nodeId string) bool {
	if s.nodeID == nodeId {
		return false
	}
	new_node := Node{nodeID: nodeId}
	//判断是否为第一次加入节点,是的话就进行一个初始化功能
	if len(s.buckets) == 0 {
		bucket := Bucket{}
		bucket.ids = append(bucket.ids, new_node.nodeID)
		s.buckets = append(s.buckets, &bucket)
		return true
	}
	result := findBucket(s.nodeID, nodeId)
	if result < 0 {
		return false
	}

	var index int
	if result >= (len(s.buckets) - 1) {
		index = len(s.buckets) - 1
		insertIntoClose(index, new_node.nodeID, s)
	} else {
		index = result
		insertIntoFar(index, new_node.nodeID, s)
	}
	return true
}

func insertIntoClose(index int, newNode string, targetNode *Node) {
	bucket := targetNode.buckets[index]
	//判断桶中是否已经存在要加入的节点
	for _, v := range bucket.ids {
		if v == newNode {
			return
		}
	}
	//判断桶是否已满，没满的话加入桶中，满的话进行扩充
	if len(bucket.ids) < 3 {
		bucket.ids = append(bucket.ids, newNode)
	} else {
		//如果桶的数量大于160个的话就不会进行分裂
		if len(targetNode.buckets) >= 160 {
			return
		}
		bucket_far := Bucket{}
		bucket_near := Bucket{}

		//将所有节点放到一起
		bucket.ids = append(bucket.ids, newNode)
		//将每个节点的距离计算出来并加入数组之中
		var distance []*big.Int
		for _, v := range bucket.ids {
			// num1, _ := strconv.ParseInt(v.nodeID, 2, 0)
			// num2, _ := strconv.ParseInt(target_node.nodeID, 2, 0)
			num1 := new(big.Int)
			num2 := new(big.Int)
			num1.SetString(v, 2)
			num2.SetString(targetNode.nodeID, 2)
			xor := new(big.Int)
			xor.Xor(num1, num2)
			distance = append(distance, xor)
		}
		//对距离数组进行筛选，选出最近的1个节点加入最远桶
		temp := distance[0]
		index_max := 0
		for i := 1; i < len(distance); i++ {
			if distance[i].Cmp(temp) > 0 {
				temp = distance[i]
				index_max = i
			}
		}

		//将最远的节点加入最远桶，最近的加入最近桶
		bucket_far.ids = append(bucket_far.ids, bucket.ids[index_max])
		for i, _ := range distance {
			if i == index_max {
				continue
			}
			bucket_near.ids = append(bucket_near.ids, bucket.ids[i])
		}

		//将bucket进行更新
		targetNode.buckets[index] = &bucket_far
		targetNode.buckets = append(targetNode.buckets, &bucket_near)
	}
}

func insertIntoFar(index int, newNode string, targetNode *Node) {
	bucket := targetNode.buckets[index]
	//查看桶中是否已经存在新的节点
	for _, v := range bucket.ids {
		if v == newNode {
			return
		}
	}
	//小于三个就更新，满了就不管（简化），事实上要进行心跳监测
	if len(bucket.ids) < 3 {
		bucket.ids = append(bucket.ids, newNode)
	}
}

// 寻找属于第几个桶
func findBucket(selfId, targetId string) int {
	num1 := new(big.Int)
	num2 := new(big.Int)
	num1.SetString(selfId, 2)
	num2.SetString(targetId, 2)

	result := new(big.Int)
	result.Xor(num1, num2)
	return (160 - len(fmt.Sprintf("%b", result)))
}

// 打印桶中的id
func (s *Bucket) printBucketContents() {
	for _, v := range s.ids {
		fmt.Printf("nodeID = %s \n", v)
	}
}

// 查看是否有相同的元素
func isDuplicate(binaryStr string, binaryStrs []string) bool {
	// 将二进制字符串转换为大整数类型
	num := new(big.Int)
	num.SetString(binaryStr, 2)

	// 判断这个大整数是否已经存在
	for _, str := range binaryStrs {
		n := new(big.Int)
		n.SetString(str, 2)
		if n.Cmp(num) == 0 {
			return true
		}
	}

	return false
}

func testInsert() {
	var binaryStrs []string
	for len(binaryStrs) < 101 {
		max := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 160), big.NewInt(1))
		// 生成一个160位的随机二进制字符串
		num, _ := rand.Int(rand.Reader, max)
		binaryStr := fmt.Sprintf("%0160b", num)

		// 检查这个二进制字符串是否已经存在
		if !isDuplicate(binaryStr, binaryStrs) {
			binaryStrs = append(binaryStrs, binaryStr)
		}
	}

	node := Node{nodeID: binaryStrs[0]}
	fmt.Println("nodeID = ", node.nodeID)
	for i, v := range binaryStrs {
		if i == 0 {
			continue
		}
		node.InsertNode(v)
	}

	for i, v := range node.buckets {
		fmt.Printf("buckets num is = %d \n", i)
		v.printBucketContents()
		fmt.Println("--------------------------")
	}
}

func testFindNode() {
	var binaryStrs []string
	for len(binaryStrs) < 206 {
		max := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 160), big.NewInt(1))
		// 生成一个160位的随机二进制字符串
		num, _ := rand.Int(rand.Reader, max)
		binaryStr := fmt.Sprintf("%0160b", num)

		// 检查这个二进制字符串是否已经存在
		if !isDuplicate(binaryStr, binaryStrs) {
			binaryStrs = append(binaryStrs, binaryStr)
		}
	}
	//介绍人节点
	var nodes []*Node
	//所有的节点
	var allNodes []*Node

	node0 := Node{nodeID: binaryStrs[0]}
	nodesMap[node0.nodeID] = &node0
	node1 := Node{nodeID: binaryStrs[1]}
	nodesMap[node1.nodeID] = &node1
	node2 := Node{nodeID: binaryStrs[2]}
	nodesMap[node2.nodeID] = &node2
	node3 := Node{nodeID: binaryStrs[3]}
	nodesMap[node3.nodeID] = &node3
	node4 := Node{nodeID: binaryStrs[4]}
	nodesMap[node4.nodeID] = &node4

	nodes = append(nodes, &node0, &node1, &node2, &node3, &node4)
	allNodes = append(allNodes, nodes...)

	//构建介绍人节点的桶
	for i := 0; i < len(nodes); i++ {
		for j := 0; j < 5; j++ {
			if i == j {
				continue
			}
			nodes[i].InsertNode(binaryStrs[j])
		}
	}

	//将200个新结点进行初始化加入网络
	for i := 5; i < len(binaryStrs); i++ {
		newNode := Node{nodeID: binaryStrs[i]}
		nodesMap[newNode.nodeID] = &newNode
		//随机选取一个介绍人节点
		num, err := rand.Int(rand.Reader, big.NewInt(5))
		if err != nil {
			// 处理错误
			return
		}

		newNode.InsertNode(nodes[num.Int64()].nodeID)
		byte1 := []byte(binaryStrs[i])
		for i, v := range byte1 {
			if v == '0' {
				byte1[i] = '1'
			} else {
				byte1[i] = '0'
			}
		}
		var tempArray []string
		far_node := Node{nodeID: string(byte1)}
		tempArray = append(tempArray, far_node.nodeID, far_node.nodeID)
		nodeIds := nodes[num.Int64()].FindNode(newNode.nodeID, tempArray)
		for _, v := range nodeIds {
			newNode.InsertNode(v)
		}
		allNodes = append(allNodes, &newNode)
	}

	for _, v := range allNodes {
		fmt.Println("Node = ", v.nodeID)
		for k, j := range v.buckets {
			fmt.Printf("buckets num is = %d \n", k)
			j.printBucketContents()
			fmt.Println("-----------------------------------------------------------------")
		}
		fmt.Println("********************************************************************************")
	}
}

func (s *Node) SetValue(key string, value []byte) bool {
	hash := sha1.Sum(value)
	hash_str := hex.EncodeToString(hash[:])
	if key != hash_str {
		return false
	}
	if s.keys[key] != nil {
		return true
	}

	//将内容存入自己的节点中
	s.keys[key] = value

	//获取到最近的桶
	keyInt := new(big.Int)
	str, _ := keyInt.SetString(key, 16)
	keyBinary := fmt.Sprintf("%160b", str)
	result := findBucket(s.nodeID, keyBinary)
	var bucket *Bucket
	if result >= (len(s.buckets) - 1) {
		bucket = s.buckets[len(s.buckets)-1]
	} else {
		bucket = s.buckets[result]
	}

	//查看找到的新的节点是否为最近的节点，如果是就进行递归，不是的话就停止递归
	index1, index2 := checkLen(len(bucket.ids))
	var flag1, flag2 bool
	if index2 != -1 {
		var maxStr string
		minStr := compareGetMin(keyBinary, bucket.ids[index1], bucket.ids[index2])
		if minStr == bucket.ids[index1] {
			maxStr = bucket.ids[index2]
		} else {
			maxStr = bucket.ids[index1]
		}
		updateIndex := isUpdated(keyBinary, SetNodes, minStr)
		if updateIndex == -1 {
			return false
		} else {
			SetNodes[updateIndex] = minStr
			flag1 = nodesMap[minStr].SetValue(key, value)
		}
		updateIndexMax := isUpdated(keyBinary, SetNodes, maxStr)
		if updateIndexMax != -1 {
			SetNodes[updateIndexMax] = maxStr
			flag2 = nodesMap[maxStr].SetValue(key, value)
		}
		if flag1 || flag2 {
			return true
		} else {
			return false
		}
	} else {
		var flag bool
		updateIndex := isUpdated(keyBinary, SetNodes, bucket.ids[0])
		if updateIndex == -1 {
			return false
		} else {
			SetNodes[updateIndex] = bucket.ids[0]
			flag = nodesMap[bucket.ids[0]].SetValue(key, value)
		}
		if flag {
			return true
		} else {
			return false
		}
	}

}

func (s *Node) GetValue(key string) []byte {
	if s.keys[key] != nil {
		hash := sha1.Sum(s.keys[key])
		hash_str := hex.EncodeToString(hash[:])
		if key != hash_str {
			return nil
		}
		return s.keys[key]
	}

	//获取到最近的桶
	keyInt := new(big.Int)
	str, _ := keyInt.SetString(key, 16)
	keyBinary := fmt.Sprintf("%160b", str)
	result := findBucket(s.nodeID, keyBinary)
	var bucket *Bucket
	if result >= (len(s.buckets) - 1) {
		bucket = s.buckets[len(s.buckets)-1]
	} else {
		bucket = s.buckets[result]
	}

	index1, index2 := checkLen(len(bucket.ids))
	if index2 != -1 {
		var maxStr string
		var value1 []byte
		var value2 []byte
		minStr := compareGetMin(keyBinary, bucket.ids[index1], bucket.ids[index2])
		if minStr == bucket.ids[index1] {
			maxStr = bucket.ids[index2]
		} else {
			maxStr = bucket.ids[index1]
		}
		updateIndex := isUpdated(keyBinary, GetNodes, minStr)
		if updateIndex == -1 {
			return nil
		} else {
			GetNodes[updateIndex] = minStr
			value1 = nodesMap[minStr].GetValue(key)
		}

		updateIndexMax := isUpdated(keyBinary, GetNodes, maxStr)
		if updateIndexMax != -1 {
			GetNodes[updateIndexMax] = maxStr
			value2 = nodesMap[maxStr].GetValue(key)
		}

		if value1 == nil && value2 == nil {
			return nil
		} else if value1 != nil && value2 != nil {
			hash1 := sha1.Sum(value1)
			hash1_str := hex.EncodeToString(hash1[:])
			hash2 := sha1.Sum(value2)
			hash2_str := hex.EncodeToString(hash2[:])
			if key != hash1_str && key != hash2_str {
				return nil
			} else {
				if key != string(hash1[:]) {
					return value2
				}
				return value1
			}
		} else {
			var hash_str string
			var value []byte
			if value1 == nil {
				hash := sha1.Sum(value2)
				hash_str = hex.EncodeToString(hash[:])
				value = value2
			} else {
				hash := sha1.Sum(value1)
				hash_str = hex.EncodeToString(hash[:])
				value = value1
			}
			if key != hash_str {
				return nil
			}
			return value
		}
	} else {
		updateIndex := isUpdated(key, GetNodes, bucket.ids[0])
		if updateIndex == -1 {
			return nil
		}
		GetNodes[updateIndex] = bucket.ids[0]
		value := nodesMap[bucket.ids[0]].GetValue(key)
		if value == nil {
			return nil
		} else {
			hash := sha1.Sum(value)
			hash_str := hex.EncodeToString(hash[:])
			if key != hash_str {
				return nil
			}
			return value
		}
	}
}

func checkLen(len int) (int, int) {
	if len > 2 {
		return GetRandom2()
	} else if len == 2 {
		return 0, 1
	} else {
		return 0, -1
	}
}

func isUpdated(targetValue string, nodes []string, compare string) int {
	targetBinary := new(big.Int)
	targetBinary.SetString(targetValue, 2)
	compareBinary := new(big.Int)
	compareBinary.SetString(compare, 2)
	minValue := compareGetMin(targetValue, nodes[0], nodes[1])
	if minValue == nodes[0] {
		maxValueBinary := new(big.Int)
		maxValueBinary.SetString(nodes[1], 2)
		resultMaxValue := new(big.Int)
		resultMaxValue.Xor(targetBinary, maxValueBinary)
		resultCom := new(big.Int)
		resultCom.Xor(targetBinary, compareBinary)
		if resultCom.Cmp(resultMaxValue) < 0 {
			return 1
		}
	} else {
		maxValueBinary := new(big.Int)
		maxValueBinary.SetString(nodes[0], 2)
		resultMaxValue := new(big.Int)
		resultMaxValue.Xor(targetBinary, maxValueBinary)
		resultCom := new(big.Int)
		resultCom.Xor(targetBinary, compareBinary)
		if resultCom.Cmp(resultMaxValue) < 0 {
			return 0
		}
	}
	return -1
}

func inverse(value string) string {
	byteArray := []byte(value)
	for i, v := range byteArray {
		if v == '0' {
			byteArray[i] = '1'
		} else {
			byteArray[i] = '0'
		}
	}
	return string(byteArray)
}

func testValue() {
	//生成100个节点，并完成网络的构建
	var binaryStrs []string
	for len(binaryStrs) < 100 {
		max := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 160), big.NewInt(1))
		// 生成一个160位的随机二进制字符串
		num, _ := rand.Int(rand.Reader, max)
		binaryStr := fmt.Sprintf("%0160b", num)

		// 检查这个二进制字符串是否已经存在
		if !isDuplicate(binaryStr, binaryStrs) {
			binaryStrs = append(binaryStrs, binaryStr)
		}
	}

	var nodes []*Node
	//初始化
	for _, v := range binaryStrs {
		node := Node{nodeID: v}
		nodes = append(nodes, &node)
		nodesMap[v] = &node
		node.keys = make(map[string][]byte)
	}
	for i, v := range nodes {
		for j, k := range binaryStrs {
			if i == j {
				continue
			}
			v.InsertNode(k)
		}
	}

	//生成200个随机字符串,并生成对应的hash
	var strs []string
	var hashs []string
	for i := 0; i < 200; i++ {
		length := 8 // 随机生成长度为8的字符串
		bytes := make([]byte, length)
		rand.Read(bytes) // 从随机源中读取指定长度的随机字节序列
		str := base64.URLEncoding.EncodeToString(bytes)
		strs = append(strs, str)
		hash := sha1.Sum([]byte(str))
		hash_str := hex.EncodeToString(hash[:])
		hashs = append(hashs, hash_str)
	}

	//寻找一个随机节点
	num, _ := rand.Int(rand.Reader, big.NewInt(100))
	for i, v := range strs {
		hashInverse := inverse(hashs[i])
		if len(SetNodes) == 0 {
			SetNodes = append(SetNodes, hashInverse, hashInverse)
		} else {
			SetNodes[0] = hashInverse
			SetNodes[1] = hashInverse
		}
		nodes[num.Int64()].SetValue(hashs[i], []byte(v))
	}

	//生成100个随机数
	var nums []int
	var isExist [200]bool
	for len(nums) < 100 {
		num1, _ := rand.Int(rand.Reader, big.NewInt(200))
		if !isExist[num1.Int64()] {
			nums = append(nums, int(num1.Int64()))
			isExist[num1.Int64()] = true
		}
	}
	for _, v := range nums {
		num2, _ := rand.Int(rand.Reader, big.NewInt(100))
		hashInverse := inverse(hashs[v])
		if len(GetNodes) == 0 {
			GetNodes = append(GetNodes, hashInverse, hashInverse)
		} else {
			GetNodes[0] = hashInverse
			GetNodes[1] = hashInverse
		}

		value := nodesMap[nodes[num2.Int64()].nodeID].GetValue(hashs[v])
		fmt.Println("key is", hashs[v])
		if value != nil {
			fmt.Println("value is ", string(value))
		} else {
			fmt.Println("Can't find value")
		}
		fmt.Println("----------------------------------")
	}

}

func main() {
	nodesMap = make(map[string]*Node)

	testValue()
}
