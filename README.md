1.获取得一组文件。文件命名为eg:W-2_1.gjf.gjf.gjf.log
  文件的序号，前面可以是'-'或者'_',后面是'.',上面取到数字1.
 1  C    Isotropic =    54.5246   Anisotropy =   125.2428
   其中取到 1，C 和54.5246，组成 C1 54.5246
后面还有
36  H    Isotropic =    29.1199   Anisotropy =     7.3113
  取到 36,H和29.1199组成H36 29.1199

组成data1.xlsx里面sheet1的前几列。

每一个个文件里面获取HF=-1387.7316846\  中获取-1387.7316846的数字，放入data.xlsx的sheet2中。

3.其中的计算公式详见data.xlsx

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o W2-CH-3ge.exe cmd/main.go