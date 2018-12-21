## Создание токена ERC20 в Ethereum для ICO

http://echain.ru/obuchenie/uroki-solidity/sozdanie-tokena-erc20-v-ethereum-dlya-ico
https://habr.com/post/344578/
https://habr.com/post/312008/

В этом уроке речь пойдет о коде токена ERC20 в сети Ethereum. Такие токены используются для ICO.
Будет разобран код токена, взятый на сайте https://ethereum.org/token:


```solidity
pragma solidity ^0.4.16;

interface tokenRecipient { function receiveApproval(address _from, uint256 _value, address _token, bytes _extraData) public; }

contract TokenERC20 {
    // Public variables of the token
    string public name;
    string public symbol;
    uint8 public decimals = 18;
    // 18 decimals is the strongly suggested default, avoid changing it
    uint256 public totalSupply;

    // This creates an array with all balances
    mapping (address => uint256) public balanceOf;
    mapping (address => mapping (address => uint256)) public allowance;

    // This generates a public event on the blockchain that will notify clients
    event Transfer(address indexed from, address indexed to, uint256 value);

    // This notifies clients about the amount burnt
    event Burn(address indexed from, uint256 value);

    /**
     * Constructor function
     *
     * Initializes contract with initial supply tokens to the creator of the contract
     */
    function TokenERC20(
        uint256 initialSupply,
        string tokenName,
        string tokenSymbol
    ) public {
        totalSupply = initialSupply * 10 ** uint256(decimals);  // Update total supply with the decimal amount
        balanceOf[msg.sender] = totalSupply;                // Give the creator all initial tokens
        name = tokenName;                                   // Set the name for display purposes
        symbol = tokenSymbol;                               // Set the symbol for display purposes
    }

    /**
     * Internal transfer, only can be called by this contract
     */
    function _transfer(address _from, address _to, uint _value) internal {
        // Prevent transfer to 0x0 address. Use burn() instead
        require(_to != 0x0);
        // Check if the sender has enough
        require(balanceOf[_from] >= _value);
        // Check for overflows
        require(balanceOf[_to] + _value > balanceOf[_to]);
        // Save this for an assertion in the future
        uint previousBalances = balanceOf[_from] + balanceOf[_to];
        // Subtract from the sender
        balanceOf[_from] -= _value;
        // Add the same to the recipient
        balanceOf[_to] += _value;
        Transfer(_from, _to, _value);
        // Asserts are used to use static analysis to find bugs in your code. They should never fail
        assert(balanceOf[_from] + balanceOf[_to] == previousBalances);
    }

    /**
     * Transfer tokens
     *
     * Send `_value` tokens to `_to` from your account
     *
     * @param _to The address of the recipient
     * @param _value the amount to send
     */
    function transfer(address _to, uint256 _value) public {
        _transfer(msg.sender, _to, _value);
    }

    /**
     * Transfer tokens from other address
     *
     * Send `_value` tokens to `_to` on behalf of `_from`
     *
     * @param _from The address of the sender
     * @param _to The address of the recipient
     * @param _value the amount to send
     */
    function transferFrom(address _from, address _to, uint256 _value) public returns (bool success) {
        require(_value <= allowance[_from][msg.sender]);     // Check allowance
        allowance[_from][msg.sender] -= _value;
        _transfer(_from, _to, _value);
        return true;
    }

    /**
     * Set allowance for other address
     *
     * Allows `_spender` to spend no more than `_value` tokens on your behalf
     *
     * @param _spender The address authorized to spend
     * @param _value the max amount they can spend
     */
    function approve(address _spender, uint256 _value) public
        returns (bool success) {
        allowance[msg.sender][_spender] = _value;
        return true;
    }

    /**
     * Set allowance for other address and notify
     *
     * Allows `_spender` to spend no more than `_value` tokens on your behalf, and then ping the contract about it
     *
     * @param _spender The address authorized to spend
     * @param _value the max amount they can spend
     * @param _extraData some extra information to send to the approved contract
     */
    function approveAndCall(address _spender, uint256 _value, bytes _extraData)
        public
        returns (bool success) {
        tokenRecipient spender = tokenRecipient(_spender);
        if (approve(_spender, _value)) {
            spender.receiveApproval(msg.sender, _value, this, _extraData);
            return true;
        }
    }

    /**
     * Destroy tokens
     *
     * Remove `_value` tokens from the system irreversibly
     *
     * @param _value the amount of money to burn
     */
    function burn(uint256 _value) public returns (bool success) {
        require(balanceOf[msg.sender] >= _value);   // Check if the sender has enough
        balanceOf[msg.sender] -= _value;            // Subtract from the sender
        totalSupply -= _value;                      // Updates totalSupply
        Burn(msg.sender, _value);
        return true;
    }

    /**
     * Destroy tokens from other account
     *
     * Remove `_value` tokens from the system irreversibly on behalf of `_from`.
     *
     * @param _from the address of the sender
     * @param _value the amount of money to burn
     */
    function burnFrom(address _from, uint256 _value) public returns (bool success) {
        require(balanceOf[_from] >= _value);                // Check if the targeted balance is enough
        require(_value <= allowance[_from][msg.sender]);    // Check allowance
        balanceOf[_from] -= _value;                         // Subtract from the targeted balance
        allowance[_from][msg.sender] -= _value;             // Subtract from the sender's allowance
        totalSupply -= _value;                              // Update totalSupply
        Burn(_from, _value);
        return true;
    }
}
```

“pragma solidity ^0.4.16” говорит о том, что код написан на Solidity версии 0.4.16 и выше (вплоть до версии 0.5.0, но не включая ее).

“interface tokenRecipient” — это интерфейс.

“contract TokenERC20” — объявление самого контракта. Всё, что дальше заключено в фигурные скобки, является кодом данного контракта. Контракты в Solidity подобны классам в объектно-ориентированных языках программирования.

“string public name; … uint256 public totalSupply;” — здесь объявляются глобальные переменные контракта. Они будут видны во всем коде данного контракта. В Solidity при объявлении переменной должен быть указан ее тип (string, uint и т.д.). Ключевое слово public называется спецификатором области видимости. Здесь public указывает на то, что переменная видна не только в данном контракте, но и в производных контрактах. К тому же, для public глобальных переменных автоматически создается геттер-функция; благодаря этому, например, в Remix, после создания смарт контракта, можно смотреть текущие значения этих переменных. В Solidity 4 спецификатора области видимости: public, external, internal и private. Подробнее об областях видимости в Solidity читайте здесь.

“mapping …” — это один из типов данных в Solidity. Маппинг можно представить в виде таблицы:

Ключ	Значение
name	Name
surname	Surname
При объявлении маппинга указывается тип ключа и тип значения. Например, в “mapping (address => uint256) public balanceOf;” ключ должен быть типа address, а значение типа uint256.

“mapping (address => mapping (address => uint256)) public allowance;” — это так называемый двойной маппинг. Как используются маппинги будет описано ниже.

“event Transfer()” и “event Burn()” — это события. События будут использоваться в некоторых функциях. События нужны для записи информации из кода в журнал блокчейн.

“function TokenERC20” — эту функцию называют конструктором контракта. Любая функция в коде контракта, название которой совпадает с названием самого контракта, становится конструктором. У контракта может быть только 1 конструктор. Функция-конструктор вызывается только 1 раз, и вызывается она в момент выгрузки контракта в сеть блокчейн. В данном коде конструктор требует 3 параметра (uint256 initialSupply, string tokenName, string tokenSymbol). Эти параметры нужно будет заполнить при выгрузке контракта в сеть. Но можно присвоить необходимые значения сразу в конструкторе, например, так:

function TokenERC20() public {
        totalSupply = 1000 * 10 ** uint256(decimals);  // Здесь вместо initialSupply сразу указано 1000
        balanceOf[msg.sender] = totalSupply;
        name = "Token Name";                          // Здесь вместо tokenName записано "Token Name"
        symbol = "TNS";                               // Здесь вместо tokenSymbol записано "TNS"
    }
Здесь нужно обратить внимание на msg.sender. В данном случае это адрес того, кто выгрузил (задеплоил) контракт в сеть блокчейн. Далее в коде будет снова встречаться msg.sender в функциях. Но в функциях значение msg.sender будет каждый раз разным в зависимости от того, кто вызывает конкретную функцию после деплоя контракта. Public-функции можно вызывать после деплоя контракта в сеть. Если Вы смотрели урок “Создание контрактов на Ethereum“, то должны понимать о чем идет речь.

Также, здесь видно, как используется mapping balanceOf.

“function _transfer” — это внутренняя функция, так как ей задана область видимости internal. Вызывать внутреннюю функцию можно только из текущего контракта или из наследующего контракта. Функция принимает 3 аргумента: address _from, address _to, uint _value. Ключевое слово require делает проверку, после которой возможно дальнейшее выполнение кода функции. Как видно, чтобы выполнилась функция _transfer, нужно чтобы выполнились 3 условия. Далее после все require создается переменная previousBalances, значение которой равно сумме балансов отправителя и получателя. Затем баланс отправителя уменьшается на значение _value, а баланс получателя наоборот увеличивается на это же значение. Затем вызывается событие Transfer. И в конце проверка на то, что сумма балансов (previousBalances) отправителя и получателя осталась такой же, как до перевода.

“function transfer” делает только одно — вызывает функцию _transfer. Обратите внимание, что здесь msg.sender — это адрес аккаунта, который вызывает функцию transfer. И этот msg.sender может отличаться от msg.sender из конструктора. Для выполнения функции transfer у msg.sender должно быть достаточное количество токенов.

“function transferFrom” осуществляет перевод токенов с указанного адреса на адрес получателя. В объявление функции добавлен “returns (bool success)”. В Solidity, если функция что-то возвращает, то в заголовке функции нужно указать тип возвращаемого значения. Для осуществления перевода с чужого кошелька у отправителя должно быть разрешение, которое предоставляется через вызов функции approve() (будет описана ниже).

“function approve” предоставляет разрешение на перевод токенов с Вашего аккаунта. Например, у Вас есть 100 токенов. Вы можете разрешить кому-то другому перевести, например, 50 токенов. Для этого Вам нужно вызвать функцию approve(), и в аргументах указать адрес доверенного аккаунта и количество доверенных токенов. При этом в функции aprrove() адресом msg.sender будет Ваш адрес (адрес аккаунта, с которого Вы вызываете функцию approve() ). Количество доверенных токенов записывается в маппинг “allowance[msg.sender][_spender] = _value;”. Вот таким образом работает двойной маппинг. То есть вы можете доверить по 10 токенов из Ваших 100 токенов. Так, будет 10 доверенных. При этом маппинге allowance msg.sender всегда будет Ваш аккаунт, а _spender будет 10 разных аккаунтов: allowance[msg.sender][_spender_1] = 10; allowance[msg.sender][_spender_2] = 10 и т.д.

