package wordlist

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Wordlist struct {
	words    []string
	filePath string
}

func NewWordlist(filePath string) *Wordlist {
	wl := &Wordlist{
		filePath: filePath,
		words:    make([]string, 0),
	}
	
	if filePath != "" {
		wl.loadFromFile()
	} else {
		wl.loadDefault()
	}
	
	return wl
}

func (w *Wordlist) GetWords() []string {
	return w.words
}

func (w *Wordlist) loadFromFile() error {
	file, err := os.Open(w.filePath)
	if err != nil {
		return fmt.Errorf("failed to open wordlist file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word != "" && !strings.HasPrefix(word, "#") {
			w.words = append(w.words, word)
		}
	}

	return scanner.Err()
}

func (w *Wordlist) loadDefault() {
	defaultWords := []string{
		"www", "mail", "ftp", "localhost", "webmail", "smtp", "pop", "ns1", "webdisk", "ns2",
		"cpanel", "whm", "autodiscover", "autoconfig", "m", "imap", "test", "ns", "blog",
		"pop3", "dev", "www2", "admin", "forum", "news", "vpn", "ns3", "mail2", "new",
		"mysql", "old", "www1", "www3", "www4", "www5", "www6", "www7", "www8", "www9",
		"www10", "api", "api1", "api2", "api3", "api4", "api5", "api6", "api7", "api8",
		"api9", "api10", "app", "app1", "app2", "app3", "app4", "app5", "app6", "app7",
		"app8", "app9", "app10", "beta", "staging", "dev1", "dev2", "dev3", "dev4", "dev5",
		"test1", "test2", "test3", "test4", "test5", "stage", "staging1", "staging2",
		"staging3", "staging4", "staging5", "demo", "demo1", "demo2", "demo3", "demo4",
		"demo5", "alpha", "beta1", "beta2", "beta3", "beta4", "beta5", "gamma", "gamma1",
		"gamma2", "gamma3", "gamma4", "gamma5", "delta", "delta1", "delta2", "delta3",
		"delta4", "delta5", "epsilon", "epsilon1", "epsilon2", "epsilon3", "epsilon4",
		"epsilon5", "zeta", "zeta1", "zeta2", "zeta3", "zeta4", "zeta5", "eta", "eta1",
		"eta2", "eta3", "eta4", "eta5", "theta", "theta1", "theta2", "theta3", "theta4",
		"theta5", "iota", "iota1", "iota2", "iota3", "iota4", "iota5", "kappa", "kappa1",
		"kappa2", "kappa3", "kappa4", "kappa5", "lambda", "lambda1", "lambda2", "lambda3",
		"lambda4", "lambda5", "mu", "mu1", "mu2", "mu3", "mu4", "mu5", "nu", "nu1", "nu2",
		"nu3", "nu4", "nu5", "xi", "xi1", "xi2", "xi3", "xi4", "xi5", "omicron", "omicron1",
		"omicron2", "omicron3", "omicron4", "omicron5", "pi", "pi1", "pi2", "pi3", "pi4",
		"pi5", "rho", "rho1", "rho2", "rho3", "rho4", "rho5", "sigma", "sigma1", "sigma2",
		"sigma3", "sigma4", "sigma5", "tau", "tau1", "tau2", "tau3", "tau4", "tau5",
		"upsilon", "upsilon1", "upsilon2", "upsilon3", "upsilon4", "upsilon5", "phi",
		"phi1", "phi2", "phi3", "phi4", "phi5", "chi", "chi1", "chi2", "chi3", "chi4",
		"chi5", "psi", "psi1", "psi2", "psi3", "psi4", "psi5", "omega", "omega1", "omega2",
		"omega3", "omega4", "omega5", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k",
		"l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "0", "1",
		"2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16",
		"17", "18", "19", "20", "21", "22", "23", "24", "25", "26", "27", "28", "29", "30",
		"31", "32", "33", "34", "35", "36", "37", "38", "39", "40", "41", "42", "43", "44",
		"45", "46", "47", "48", "49", "50", "51", "52", "53", "54", "55", "56", "57", "58",
		"59", "60", "61", "62", "63", "64", "65", "66", "67", "68", "69", "70", "71", "72",
		"73", "74", "75", "76", "77", "78", "79", "80", "81", "82", "83", "84", "85", "86",
		"87", "88", "89", "90", "91", "92", "93", "94", "95", "96", "97", "98", "99", "100",
	}
	
	w.words = defaultWords
}

func (w *Wordlist) AddWord(word string) {
	w.words = append(w.words, word)
}

func (w *Wordlist) RemoveWord(word string) {
	for i, wordItem := range w.words {
		if wordItem == word {
			w.words = append(w.words[:i], w.words[i+1:]...)
			break
		}
	}
}

func (w *Wordlist) GetCount() int {
	return len(w.words)
}

func (w *Wordlist) SaveToFile(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create wordlist file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, word := range w.words {
		_, err := writer.WriteString(word + "\n")
		if err != nil {
			return fmt.Errorf("failed to write word to file: %v", err)
		}
	}

	return writer.Flush()
}
