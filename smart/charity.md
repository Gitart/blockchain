Рассмотрим код контракта Charity.sol по логическим блокам. Сначала идет описание всех нужных нам переменных:

    uint public debatingPeriodInMinutes; // время на голосование
    Proposal[] public proposals; // массив предложений пожертвований, структура описана далее
    uint public numProposals; // количество элементов в массиве proposals
    uint public numMembers; //количество элементов в массиве members
    mapping (address => uint) public memberId; // маппинг адреса аккаунта на индекс в массиве members
    address[] public members; // массив зарегистрированных аккаунтов

Маппинг позволяет по адресу пользователя получить его индекс в массиве. Если пользователя с таким адресом не существует, то будет возвращен индекс 0. На этом будет основана далее функция, определяющая, зарегистрирован ли пользователь с данным адресом. Но это вносит требование для массива: пользователи должны храниться в массиве начиная с индекса 1. Код, отвечающий за эту логику будет рассмотрен дальше. А пока идет структура для хранения предложения.

    struct Proposal {
        address recipient; // получатель
        uint amount; // размер пожертвования
        string description; // описание пожертвования
        uint votingDeadline; // дедлайн
        bool executed; // флаг, что пожертвование уже совершено или отклонено
        bool proposalPassed; // флаг, что голосование одобрило пожертвование
        uint numberOfVotes; // количество голосов
        int currentResult; // сумма голосов, “за“ = +1, “против” = -1
        Vote[] votes; // список голосов с адресом каждого проголосовавшего и значением голоса, структура описана далее
        mapping (address => bool) voted; // маппинг для быстрой проверки, проголосовал ли аккаунт с таким-то адресом
    }

Структура голосов, складывается в массив для каждого предложения

    struct Vote {
        bool inSupport; // значение голоса
        address voter; // и адрес аккаунта проголосовавшего
    }

Рассмотрим модификатор, который позволит нам контролировать, что выполнение методов, к которым он будет добавлен, возможно только если пользователь зарегистрирован. Код проверки, как уже было сказано, основан на том, что несуществующие элементы маппинга дают индекс 0, а пользователей будем хранить начиная с индекса 1:

    modifier onlyMembers {
        require (memberId[msg.sender] != 0); // методы будут выполнять такой код проверки
        _; // код модифицируемого метода подставляется вместо знака подчеркивания
    }

msg — это структура, по которой можно получить информацию о вызывающем. В данном случае msg.sender — это адрес аккаунта, который вызвал метод с этим модификатором. 

Опишем конструктор нашего контракта, который будет выполняться при деплое. Все, что требуется задавать — время, которое выделяется для голосования за каждое предложение. Кроме этого увеличиваем размер массива members, потому что добавлять пользователей будем исходя из размера, а нулевой элемент остается зарезервированным.

    function Charity( uint _minutesForDebate ) payable { // payable означает, что вместе с транзакцией можно отправить эфир, этот эфир зачислится на счет контракта
        debatingPeriodInMinutes = _minutesForDebate;
        members.length++;
    }

Функция для добавления пользователя: 

    function addMember(address _targetMember) {
        if (memberId[_targetMember] == 0) { // 0 является признаком, что пользователь не зарегистрирован
            uint id;
            memberId[_targetMember] = members.length; // индексом будет номер еще не добавленного элемента, сохраняем его в маппинге
            id = members.length++; // сохраняем индекс и увеличиваем размер массива
            members[id] = _targetMember; // сохраняем адрес в массиве
        }
    }

Заметьте функцию require — она пришла на замену throw в более старых версиях solidity. В require передается true или false, если это false — то срабатывает обработчик аналогичный throw — откатывается вся транзакция. 
Чтобы можно было проверить, находится ли адрес в списке пользователей используем такую функцию:

function isMember( address _targetMember ) constant returns ( bool ) {
        return ( memberId[_targetMember] != 0 );
    }

Следующая функция — для создания предложения, принимает адрес получателя пожертвования, количество эфира в wei и строку с описанием. К этой функции применяется модификатор onlyMembers, это значит до выполнения всего кода произойдет проверка, что вызывающий аккаунт зарегистрирован. Здесь вы увидите такие преобразования как 1 ether и 1 minutes. Полный список таких суффиксов можете посмотреть здесь, они сделаны для удобства и могут применяться только к значениям, но не к переменным. Но чтобы применить к переменной достаточно просто добавить 1 к суффиксу, что и сделано в нашем случае для преобразования в секунды.

    function newProposal(
            address _beneficiary, // получатель пожертвования
            uint _weiAmount, // размер пожертвования в wei
            string _description // произвольная строка
        )
            onlyMembers
            returns (uint proposalID)
    {
        require( _weiAmount <= (1 ether) ); // ограничим пожертвование в 1 эфира 
        proposalID = proposals.length++; // увеличение размера массива на 1
        Proposal storage p = proposals[proposalID]; // далее идет присвоение элементу значений
        p.recipient = _beneficiary;
        p.amount = _weiAmount;
        p.description = _description;
        p.votingDeadline = now + debatingPeriodInMinutes * 1 minutes; // расчет дедлайна по текущему времени и времени на голосование
        p.executed = false; // обнуление флагов завершенности транзакции
        p.proposalPassed = false;
        p.numberOfVotes = 0;
        numProposals = proposalID + 1; // сохранение размера массива

        return proposalID;
    }

Заметьте здесь ключевое слово now — это текущее время, но не на момент вызова транзакции, а на момент создания блока. Поэтому дедлайн будет отсчитываться с момента, когда предложение уже будет создано на блокчейне.

Несмотря на то, что proposals у нас public, получать таким образом можно только простейшие поля в виде массива. То есть вызвав в контракте метод например proposals(1), мы получим предложение с индексом 1 в виде массива { recipient, amount, description, votingDeadline, executed, proposalPassed, numberOfVotes, currentResult }, а массивы votes и voted внутри структуры не вернутся. Но нам нужна информация о том, проголосовал ли пользователь за определенное предложение, чтобы отображать его голос или дать возможность проголосовать. И желательно сделать это в одно обращение, поэтому мы получаем эту информацию когда читаем структуру Proposal для отображения в нашем приложении с помощью специальной функции getProposal, которая принимает аккаунт, для которого нужен статус голоса и идентификатор предложения.

    function getProposal( address _member, uint _proposalNumber ) constant
        returns ( address, // описываем типы в возвращаемом массиве
                  uint,
                  string,
                  uint,
                  bool,
                  bool,
                  uint,
                  int,
                  int ) {
        Proposal memory proposal = proposals[ _proposalNumber ]; // берем элемент для удобства
        int vote = getVoted( _member, _proposalNumber ); // используем вспомогательную функцию (описана позже) для получения информации о голосе конкретного пользователя
        return ( proposal.recipient,
                 proposal.amount,
                 proposal.description,
                 proposal.votingDeadline,
                 proposal.executed,
                 proposal.proposalPassed,
                 proposal.numberOfVotes,
                 proposal.currentResult,
                 vote ); // высылаем массив в соответствии с ожидаемыми типами
    }

А это вспомогательная функция, которая ищет как проголосовал конкретный пользователь в конкретном предложении. Возвращаться будет: 0 — если пользователь не проголосовал, 1 — если пользователь проголосовал “за”, -1 — если проголосовал “против”.

    function getVoted(address _member, uint _proposalNumber) constant returns(int)
    {
        Proposal storage p = proposals[_proposalNumber];
        int result = 0;
        int true_int = 1;
        int false_int = -1; // определяем возвращаемые значения
        for (uint i = 0; i < p.numberOfVotes; i++)
        {
            if (p.votes[i].voter == _member) // ищем нужного пользователя перебором
            {
                result = p.votes[i].inSupport ? true_int : false_int;
                break; // если находим выходим и возвращаем значение
            }
        }
        return result;
    }

Голосование: для предложения с конкретным номером отдаем голос true (за) или false (против).

    function vote(
            uint _proposalNumber, // номер предложения, за которое отдается голос
            bool _supportsProposal // голос
    )
            onlyMembers
            returns (uint voteID)
    {
        Proposal storage p = proposals[_proposalNumber];    // для удобства возьмем нужный элемент из массива
        require (p.voted[msg.sender] != true);      // не продолжать если пользователь уже голосовал
        p.voted[msg.sender] = true;                // отметить пользователя как проголосовавшего
        p.numberOfVotes++;                 // увеличить количество проголосовавших для предложения
        if (_supportsProposal) {                        // если проголосовали “за”
            p.currentResult++;                          // увеличить результат на 1
        } else {                                        // если против
            p.currentResult--;                          // уменьшить результат на 1
        }
        voteID = p.votes.length++; // добавление нового голоса в массив голосов
        p.votes[voteID] = Vote({inSupport: _supportsProposal, voter: msg.sender}); // инициализация структуры
        return p.numberOfVotes;
    }

И последняя функция executeProposal служит для завершения голосования и отправки (или неотправки) эфира на адрес получателя.

    function executeProposal(uint _proposalNumber) { // выполнить предложение с таким номером
        Proposal storage p = proposals[_proposalNumber];

        require ( !(now < p.votingDeadline || p.executed) ); // предложение должно 1) пройти дедлайн, 2) не быть уже выполненным
        p.executed = true; // помечаем как выполненное

        if (p.currentResult > 0) { // если большинство проголосовало “за”
            require ( p.recipient.send(p.amount) ); // отправить эфир получателю
            p.proposalPassed = true; // пометить, что предложение одобрено
        } else { // если “за” проголосовало не большинство
            p.proposalPassed = false; // пометить, что предложение отклонено и ничего не отправлять
        }
    }

В конце присутствует пустая функция с модификатором payable.

function () payable {}

Это нужно для того, чтобы на адрес контракта можно было присылать эфир. Вообще пустая функция — это функция, которая принимает и обрабатывает все сообщения, которые не являются вызовом функций. Все, что нам требуется — это сделать ее payable, тогда отправленный газ просто зачислится на контракт без каких-либо дополнительных действий. Но заметьте, что на других функциях этого модификатора нет, поэтому в нашем случае нельзя отправлять эфир например с вызовом addMember.

Вариант приложения с использованием Web3.js

Основной сценарий приложения:

Пользователь подключается к сети Ropsten через MetaMask
Если на аккаунте нет эфира, будет невозможно выполнить ни одну транзакцию. Мы добавили функцию получения эфира, которая становится доступна при балансе аккаунта меньше 0.1 эфира. Реализовано это через сторонний сервис, на который делается ajax запрос с адресом, на который нужно перевести эфир.
Основные действия со смарт контрактом доступны только после того, как пользователь станет участником организации. Для этого вызывается метод addMember в смарт контракте.
Участник организации может создать Предложение о переводе средств (далее Proposal), или проголосовать за уже существующее.
Когда истекает время для Proposal (время создания + 5 минут), появляется возможность его завершить, в результате чего, в зависимости от распределения голосов, эфир будет переведен на указанный адрес, или нет.

Демонстрация приложения доступна по ссылке — MetaMask версия.
Исходный код здесь. 

Еще раз обращаем ваше внимание на то, что текущая версия Web3.js — 0.20.1. Но уже готовится к релизу версия 1.0, в которой изменения достаточно существенны. Как мы говорили выше, MetaMask встраивает web3 в страницу, и его можно сразу использовать. Но учитывая то, что библиотека активно развивается, а нам нужно гарантировать работоспособность приложения для пользователя, необходимо использовать свою залоченную версию, и переопределять объект web3, который встраивает MetaMask. Мы делаем это здесь в следующем методе:

 initializeWeb3() {
    if (typeof web3 !== 'undefined') { // если MetaMask заинжектил библиотеку
      const defaultAccount = web3.eth.defaultAccount; // сохраняем привязанный аккаунт
      window.web3 = new Web3(web3.currentProvider); // инициализируем свою библиотеку
      window.web3.eth.defaultAccount = defaultAccount; // возвращаем привязанный аккаунт
    }  }

Делать это нужно после события window.onload.
Одна неочевидная проблема, которая решается в этом коде — если просто сделать window.web3 = new Web3(web3.currentProvider) как предлагается в официальной документации, то не подхватывается аккаунт по умолчанию. 
Еще в MetaMask, как уже писалось, можно выбирать сеть из списка. У нас используются адреса контрактов в сети Ropsten, если попытаться подключаться по этим адресам в других сетях — результат будет непредсказуем. Поэтому прежде чем предоставлять доступ к приложению, нужно проверить в той ли сети находится пользователь. Получить идентификатор сети можно с помощью команды:

web3.version.getNetwork(function (err, netId) {});

Мы делаем эту проверку здесь и сравниваем результат с id для сети Ropsten — это 3.

Список id всех сетей можно увидеть например здесь в описании net_version. 

Вся логика работы с блокчейном находится в файле blockchain.js.

Здесь есть два типа функций — функции для получения данных из блокчейна и функции изменяющие данные в блокчейне. Большинство методов из web3.js выполняются асинхронно и принимают callback в качестве последнего параметра. Поскольку зачастую приходится вызывать несколько методов для получения данных, и вызов некоторых из них зависит от результата работы других — удобно использовать промисы. В версии 1.0 web3.js асинхронные методы возвращают промисы по умолчанию.

Приведем один пример получения информации из блокчейна:
Функция getCurrentAccountInfo возвращает адрес текущего аккаунта, баланс и флаг того, является ли данный аккаунт участником организации.

Blockchain.prototype.getCurrentAccountInfo = function() {
  const address = this.address;
  if (address == undefined) {
    return Promise.resolve({});
  }

  const balancePromise = new Promise(function(resolve, reject) {
    web3.eth.getBalance(address, function(err, res) {
      err ? reject(err) : resolve(web3.fromWei(res).toNumber());
    });
  });

  const authorizedPromise = new Promise(function(resolve, reject) {
    this.contractInstance.isMember(address, function(err, res) {
      err ? reject(err) : resolve(res);
    });
  }.bind(this));

  return new Promise(function(resolve, reject) {
    Promise.all([balancePromise, authorizedPromise]).then(function(data) {
      resolve({
        address: address,
        balance: data[0],
        isMember: data[1]
      });
    });
  });
};

Рассмотрим теперь функцию изменения данных в блокчейне, например функция добавления участника организации.

Blockchain.prototype.becomeMember = function() {
  return new Promise(function(resolve, reject) {
    this.contractInstance.addMember(this.address, function(err, res) {
      err ? reject(err) : resolve(res);
    });
  }.bind(this));
};

Как видим, синтаксис ничем не отличается от предыдущего примера, вот только выполнение этой функции повлечет создание транзакции, для изменения данных в блокчейне.
При вызове любой функции смарт контракта, в результате которой создается транзакция, MetaMask предлагает пользователю подтвердить эту транзакцию или отклонить ее. Если пользователь подтверждает транзакцию, то функция возвращает хеш транзакции.
Один неочевидный момент — это как узнать выполнилась транзакция успешно или нет.
Определить статус транзакции можно на основании кол-ва газа, которое было использовано. Если использовано максимально доступное кол-во газа, то либо в ходе выполнения возникла ошибка, либо газа не хватило для выполнения транзакции. Проверку статуса мы делаем следующим образом.

Blockchain.prototype.checkTransaction = function(transaction) {
  const txPromise = new Promise(function(resolve, reject) {
    web3.eth.getTransaction(transaction.transactionHash, function(err, res) {
      err ? reject(err) : resolve(res);
    });
  });

  const txReceiptPromise = new Promise(function(resolve, reject) {
    web3.eth.getTransactionReceipt(transaction.transactionHash, function(err, res) {
      err ? reject(err) : resolve(res);
    });
  });

  return new Promise(function(resolve, reject) {
    Promise.all([txPromise, txReceiptPromise]).then(function(res) {
      const tx = res[0];
      const txReceipt = res[1];
      const succeeded = txReceipt && txReceipt.blockNumber && txReceipt.gasUsed < tx.gas;
      const failed = txReceipt && txReceipt.blockNumber && txReceipt.gasUsed == tx.gas;

      let state = transactionStates.STATE_PENDING;
      if (succeeded) {
        state = transactionStates.STATE_SUCCEEDED;
      } else if (failed) {
        state = transactionStates.STATE_FAILED;
      }

      resolve(state);
    });
  });
};
