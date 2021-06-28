package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

const(
	FUKA_Ma = 8.7
	FUKA_Wa = 7.9
	CONT_Ma = 7.9
	CONT_Wa = 7.2
	BOSS_Ma = 7.5
	BOSS_Wa = 6.6
	COLL_Ma = 8.1
	COLL_Wa = 8.2

	FUKA_Mb = 0.076
	FUKA_Wb = 0.048
	CONT_Mb = -0.089
	CONT_Wb = -0.056
	BOSS_Mb = -0.097
	BOSS_Wb = -0.097
	COLL_Mb = -0.097
	COLL_Wb = -0.097
)



func failOnError(err error) {
	if err != nil {
		log.Fatal("Error:", err)
	}
}

func main() {
	flag.Parse()

	// ログファイル準備
	logfile, err := os.OpenFile("./log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	failOnError(err)
	defer logfile.Close()

	log.SetOutput(logfile)

	//入力ファイル準備
	infile, err := os.Open(flag.Arg(0))
	failOnError(err)
	defer infile.Close()

	//書き込みファイル準備
	outfile, err := os.Create("./ストレスチェック集団分析.csv")
	failOnError(err)
	defer outfile.Close()

	reader := csv.NewReader(transform.NewReader(infile, japanese.ShiftJIS.NewDecoder()))
	writer := csv.NewWriter(transform.NewWriter(outfile, japanese.ShiftJIS.NewEncoder()))
	writer.UseCRLF = true

	log.Print("集団分析 Start\r\n")
	// タイトル行を取得
	_, err = reader.Read() // 1行読み出す
	if err != io.EOF {
		failOnError(err)
	}

	//タイトル行の書き込み
	writer.Write(recHeadStr())

	var out_record []string
	var A1,A2,B1,B2 []string
	for {
		record, err := reader.Read() // 1行読み出す
		if err == io.EOF {
			break
		} else {
			failOnError(err)
		}

		A1 = append(A1, record[196]) //仕事の量的負担
		A2 = append(A2, record[197]) //仕事のコントロール
		B1 = append(B1, record[198]) //上司の支援 
		B2 = append(B2, record[199]) //同僚の支援
	}

	//健康リスクの平均計算
	sRisk := make([]float64, 4)
	sRisk[0] = average(A1)
	sRisk[1] = average(A2)
	sRisk[2] = average(B1)
	sRisk[3] = average(B2)
	out_record = append(out_record, floatToString(sRisk)...)

	//健康リスク計算
	FukaCont_m := 100 * math.Exp((sRisk[0] - FUKA_Ma) * FUKA_Mb + (sRisk[1] - CONT_Ma) * CONT_Mb)
	Syokuba_m  := 100 * math.Exp((sRisk[2] - BOSS_Ma) * BOSS_Mb + (sRisk[3] - COLL_Ma) * COLL_Mb)
	SogoRisk_m := FukaCont_m * Syokuba_m / 100
	FukaCont_w := 100 * math.Exp((sRisk[0] - FUKA_Wa) * FUKA_Wb + (sRisk[1] - CONT_Wa) * CONT_Wb)
	Syokuba_w  := 100 * math.Exp((sRisk[2] - BOSS_Wa) * BOSS_Wb + (sRisk[3] - COLL_Wa) * COLL_Wb)
	SogoRisk_w := FukaCont_w * Syokuba_w / 100



	writer.Write(out_record)

	writer.Flush()
	log.Print("集団分析 Finiseh !\r\n")

}

func recHeadStr() []string {
	var Head []string
	Head = append(Head, "仕事の量的負担")
	Head = append(Head, "仕事のコントロール")
	Head = append(Head, "上司の支援")
	Head = append(Head, "同僚の支援")
	Head = append(Head, "量・コントロール_男性")
	Head = append(Head, "職場の支援_男性")
	Head = append(Head, "総合健康リスク_男性")
	Head = append(Head, "量・コントロール_女性")
	Head = append(Head, "職場の支援_女性")
	Head = append(Head, "総合健康リスク_女性")
	Head = append(Head, "量的負担_平均")
	Head = append(Head, "質的負担_平均")
	Head = append(Head, "身体負担_平均")
	Head = append(Head, "対人関係_平均")
	Head = append(Head, "職場環境_平均")
	Head = append(Head, "コントロール_平均")
	Head = append(Head, "技能活用_平均")
	Head = append(Head, "適性度_平均")
	Head = append(Head, "働き甲斐_平均")
	Head = append(Head, "活気_平均")
	Head = append(Head, "いらいら感_平均")
	Head = append(Head, "疲労感_平均")
	Head = append(Head, "不安感_平均")
	Head = append(Head, "抑うつ感_平均")
	Head = append(Head, "身体愁訴_平均")
	Head = append(Head, "上司支援_平均")
	Head = append(Head, "同僚支援_平均")
	Head = append(Head, "家族・友人支援_平均")
	Head = append(Head, "満足度_平均")
	return Head

}

func round(f float64, d int) float64 {
	//四捨五入の計算 dは小数点の位置
	shift := math.Pow(10, float64(d))
	return math.Floor(f*shift+.5) / shift
}

func floatToString(f []float64) []string {
	s := make([]string, len(f))
	for n := range f {
		s[n] = fmt.Sprint(f[n])
	}
	return s
}

func average(s []string) float64 {
	var sum int
	for _, v := range s {
		vint, _ := strconv.Atoi(v)
		sum = sum + vint
	}
	return round(float64(sum) / float64(len(s)), 1)
}

func sRiskCalman(f []float64) float64 {

